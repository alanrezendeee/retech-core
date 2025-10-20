package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Migration representa uma migração/seed
type Migration struct {
	Version     string
	Description string
	Apply       func(ctx context.Context, db *mongo.Database, log zerolog.Logger) error
}

// MigrationRecord registra migrations executadas
type MigrationRecord struct {
	Version     string    `bson:"version"`
	Description string    `bson:"description"`
	AppliedAt   time.Time `bson:"appliedAt"`
}

// MigrationManager gerencia as migrations
type MigrationManager struct {
	db         *mongo.Database
	log        zerolog.Logger
	migrations []Migration
}

// NewMigrationManager cria um novo gerenciador de migrations
func NewMigrationManager(db *mongo.Database, log zerolog.Logger) *MigrationManager {
	return &MigrationManager{
		db:  db,
		log: log,
		migrations: []Migration{
			{
				Version:     "001_seed_estados",
				Description: "Popular estados brasileiros",
				Apply:       seedEstados,
			},
			{
				Version:     "002_seed_municipios",
				Description: "Popular municípios brasileiros",
				Apply:       seedMunicipios,
			},
		},
	}
}

// Run executa as migrations pendentes
func (m *MigrationManager) Run(ctx context.Context) error {
	coll := m.db.Collection("migrations")

	for _, migration := range m.migrations {
		// Verifica se já foi aplicada
		count, err := coll.CountDocuments(ctx, bson.M{"version": migration.Version})
		if err != nil {
			return fmt.Errorf("erro ao verificar migration %s: %w", migration.Version, err)
		}

		if count > 0 {
			m.log.Info().Msgf("[migration] %s já aplicada, pulando", migration.Version)
			continue
		}

		// Aplica a migration
		m.log.Info().Msgf("[migration] Aplicando %s: %s", migration.Version, migration.Description)
		start := time.Now()

		if err := migration.Apply(ctx, m.db, m.log); err != nil {
			return fmt.Errorf("erro ao aplicar migration %s: %w", migration.Version, err)
		}

		// Registra como aplicada
		record := MigrationRecord{
			Version:     migration.Version,
			Description: migration.Description,
			AppliedAt:   time.Now(),
		}
		if _, err := coll.InsertOne(ctx, record); err != nil {
			return fmt.Errorf("erro ao registrar migration %s: %w", migration.Version, err)
		}

		m.log.Info().Msgf("[migration] %s aplicada com sucesso em %v", migration.Version, time.Since(start))
	}

	return nil
}

// seedEstados popula os estados
func seedEstados(ctx context.Context, db *mongo.Database, log zerolog.Logger) error {
	repo := storage.NewEstadosRepo(db)

	// Verifica se já existem dados
	count, err := repo.Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Info().Msgf("[seed] Estados já populados (%d registros), pulando", count)
		return nil
	}

	// Procura o arquivo estados.json
	seedFile := findSeedFile("estados.json")
	if seedFile == "" {
		return fmt.Errorf("arquivo estados.json não encontrado")
	}

	log.Info().Msgf("[seed] Carregando estados de: %s", seedFile)

	// Lê o arquivo
	data, err := os.ReadFile(seedFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo estados.json: %w", err)
	}

	var estados []domain.Estado
	if err := json.Unmarshal(data, &estados); err != nil {
		return fmt.Errorf("erro ao fazer parse de estados.json: %w", err)
	}

	// Insere no banco
	if err := repo.InsertMany(ctx, estados); err != nil {
		return fmt.Errorf("erro ao inserir estados: %w", err)
	}

	log.Info().Msgf("[seed] %d estados inseridos com sucesso", len(estados))
	return nil
}

// seedMunicipios popula os municípios
func seedMunicipios(ctx context.Context, db *mongo.Database, log zerolog.Logger) error {
	repo := storage.NewMunicipiosRepo(db)

	// Verifica se já existem dados
	count, err := repo.Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Info().Msgf("[seed] Municípios já populados (%d registros), pulando", count)
		return nil
	}

	// Procura o arquivo municipios.json
	seedFile := findSeedFile("municipios.json")
	if seedFile == "" {
		return fmt.Errorf("arquivo municipios.json não encontrado")
	}

	log.Info().Msgf("[seed] Carregando municípios de: %s", seedFile)

	// Lê o arquivo
	data, err := os.ReadFile(seedFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo municipios.json: %w", err)
	}

	var municipios []domain.Municipio
	if err := json.Unmarshal(data, &municipios); err != nil {
		return fmt.Errorf("erro ao fazer parse de municipios.json: %w", err)
	}

	log.Info().Msgf("[seed] Inserindo %d municípios (isso pode demorar)...", len(municipios))

	// Insere no banco em lotes
	if err := repo.InsertMany(ctx, municipios); err != nil {
		return fmt.Errorf("erro ao inserir municípios: %w", err)
	}

	log.Info().Msgf("[seed] %d municípios inseridos com sucesso", len(municipios))
	return nil
}

// findSeedFile procura o arquivo de seed em diversos locais
func findSeedFile(filename string) string {
	// Possíveis localizações (em ordem de prioridade)
	locations := []string{
		// 1. Diretório seeds (padrão - funciona local e Docker)
		filepath.Join("seeds", filename),
		// 2. Diretório /app/seeds (Docker com volume montado)
		filepath.Join("/app", "seeds", filename),
		// 3. Diretório atual
		filename,
		// 4. Downloads do usuário (desenvolvimento local)
		filepath.Join(os.Getenv("HOME"), "Downloads", filename),
		// 5. Diretório data
		filepath.Join("data", filename),
		// 6. Caminho absoluto no workdir
		filepath.Join(".", "seeds", filename),
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	return ""
}

