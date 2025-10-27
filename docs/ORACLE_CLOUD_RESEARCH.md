# 🔍 PESQUISA: Oracle Cloud Infrastructure (OCI) - Automação Completa

## 📋 **OBJETIVO:**
Criar scripts de automação `.sh` para provisionar e gerenciar infraestrutura na Oracle Cloud, similar ao que Railway faz via interface gráfica, mas 100% via CLI.

---

## 🎯 **O QUE QUEREMOS AUTOMATIZAR:**

### **1. Provisionamento Inicial**
- ✅ Criar VM nova
- ✅ Configurar rede (VCN, subnet, security list)
- ✅ Instalar Docker + Docker Compose
- ✅ Configurar firewall
- ✅ Setup de usuários e permissões

### **2. Deploy de Serviços**
- ✅ Subir Redis (via Docker Hub)
- ✅ Subir PostgreSQL/MongoDB (via Docker Hub)
- ✅ Subir aplicação backend
- ✅ Configurar variáveis de ambiente
- ✅ Setup de volumes persistentes

### **3. CI/CD**
- ✅ Integração com GitHub
- ✅ Deploy automático por branch
- ✅ Rollback automático em caso de erro
- ✅ Webhooks para notificações

### **4. Monitoramento & Logs**
- ✅ Configurar logs centralizados
- ✅ Métricas de CPU, RAM, Disco
- ✅ Alertas de saúde da aplicação
- ✅ Backup automático

### **5. Escalabilidade**
- ✅ Aumentar vCPU dinamicamente
- ✅ Aumentar memória
- ✅ Aumentar storage
- ✅ Load balancer (se necessário)

---

## 🛠️ **OCI CLI - PRINCIPAIS COMANDOS:**

### **Instalação**
```bash
bash -c "$(curl -L https://raw.githubusercontent.com/oracle/oci-cli/master/scripts/install/install.sh)"
```

### **Configuração**
```bash
oci setup config
# Pede:
# - User OCID
# - Tenancy OCID
# - Region (sa-saopaulo-1 para São Paulo)
# - API Key
```

### **Criar Compute Instance**
```bash
oci compute instance launch \
  --availability-domain "xxx:SA-SAOPAULO-1-AD-1" \
  --compartment-id "ocid1.compartment.oc1.xxx" \
  --shape "VM.Standard.E2.1.Micro" \
  --display-name "retech-core-prod" \
  --image-id "ocid1.image.oc1.xxx" \
  --subnet-id "ocid1.subnet.oc1.xxx" \
  --assign-public-ip true \
  --ssh-authorized-keys-file ~/.ssh/id_rsa.pub \
  --user-data-file cloud-init.sh
```

### **Listar Instâncias**
```bash
oci compute instance list \
  --compartment-id "ocid1.compartment.oc1.xxx" \
  --availability-domain "xxx:SA-SAOPAULO-1-AD-1"
```

### **Iniciar/Parar Instância**
```bash
# Iniciar
oci compute instance action --action START --instance-id "ocid1.instance.oc1.xxx"

# Parar
oci compute instance action --action STOP --instance-id "ocid1.instance.oc1.xxx"
```

### **Atualizar Recursos (Shape)**
```bash
oci compute instance update \
  --instance-id "ocid1.instance.oc1.xxx" \
  --shape "VM.Standard.E2.2"
```

### **Criar Volume (Storage)**
```bash
oci bv volume create \
  --availability-domain "xxx:SA-SAOPAULO-1-AD-1" \
  --compartment-id "ocid1.compartment.oc1.xxx" \
  --display-name "retech-storage" \
  --size-in-gbs 100
```

### **Anexar Volume**
```bash
oci compute volume-attachment attach \
  --instance-id "ocid1.instance.oc1.xxx" \
  --type paravirtualized \
  --volume-id "ocid1.volume.oc1.xxx"
```

---

## 💰 **ORACLE ALWAYS FREE TIER - RECURSOS GRATUITOS:**

### **Compute (VMs)**
```
✅ 2x VMs AMD:
   - Shape: VM.Standard.E2.1.Micro
   - vCPU: 1/8 OCPU (equivalente a 1 core compartilhado)
   - RAM: 1 GB
   - Network: 1 Gbps

OU

✅ 4x VMs ARM Ampere A1:
   - vCPU: 4 cores no total (distribuível)
   - RAM: 24 GB no total (distribuível)
   - Network: 1 Gbps
   - Exemplo: 1 VM com 4 cores + 24GB ou 4 VMs com 1 core + 6GB cada
```

### **Storage**
```
✅ Block Volume: 200 GB total
✅ Object Storage: 20 GB
✅ Archive Storage: 10 GB
```

### **Network**
```
✅ Outbound Data Transfer: 10 TB/mês
✅ Load Balancer: 1 instância (10 Mbps)
✅ VCN: Ilimitado
✅ Public IPs: 2 IPv4 reservados
```

### **Database (Autonomous)**
```
✅ 2x Autonomous Databases (ATP ou ADW)
   - Storage: 20 GB cada
   - OCPU: 1 cada
```

### **Outros Serviços**
```
✅ Monitoring: 500M ingestion datapoints/mês
✅ Notifications: 1M por mês
✅ Logging: 10 GB/mês
```

---

## 🐳 **DOCKER + DOCKER COMPOSE NO OCI:**

### **Cloud-Init Script (executa na criação da VM)**
```bash
#!/bin/bash

# Atualizar sistema
apt-get update && apt-get upgrade -y

# Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Instalar Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Criar usuário para deploy
useradd -m -s /bin/bash deploy
usermod -aG docker deploy

# Criar diretórios
mkdir -p /app
mkdir -p /app/logs
mkdir -p /app/data
chown -R deploy:deploy /app

# Instalar ferramentas úteis
apt-get install -y git curl wget vim htop

echo "✅ Setup completo!"
```

---

## 🔐 **VARIÁVEIS DE AMBIENTE:**

### **Estratégia 1: Arquivo .env no servidor**
```bash
# /app/.env
NODE_ENV=production
PORT=8080
MONGO_URI=mongodb://mongo:27017/retech
REDIS_URL=redis://redis:6379
JWT_SECRET=xxx
APIKEY_HASH_SECRET=xxx
```

### **Estratégia 2: Secrets via OCI Vault (SEGURO)**
```bash
# Criar secret
oci vault secret create-base64 \
  --compartment-id "xxx" \
  --secret-name "retech-api-key" \
  --vault-id "xxx" \
  --key-id "xxx" \
  --secret-content-content "base64encodedvalue"

# Recuperar secret
oci secrets secret-bundle get \
  --secret-id "xxx" \
  --query 'data."secret-bundle-content".content' \
  --raw-output | base64 -d
```

### **Estratégia 3: Instance Metadata (menos seguro)**
```bash
# Definir metadata na criação
oci compute instance launch \
  --metadata '{"env_vars": "{\"PORT\":\"8080\",\"NODE_ENV\":\"production\"}"}'

# Ler metadata dentro da VM
curl -H "Authorization: Bearer Oracle" http://169.254.169.254/opc/v2/instance/metadata/
```

---

## 📊 **LOGS & MONITORAMENTO:**

### **Logs Centralizados**
```bash
# Instalar OCI Logging Agent
wget https://objectstorage.sa-saopaulo-1.oraclecloud.com/n/xxx/b/xxx/o/install.sh
bash install.sh

# Configurar log sources
cat > /etc/oci-logging-agent/config.json << EOF
{
  "sources": [
    {
      "name": "app-logs",
      "type": "file",
      "path": "/app/logs/*.log"
    },
    {
      "name": "docker-logs",
      "type": "docker",
      "containers": ["retech-api", "mongo", "redis"]
    }
  ],
  "destination": {
    "log-object-id": "ocid1.log.oc1.xxx"
  }
}
EOF
```

### **Métricas (CPU, RAM, Disco)**
```bash
# Habilitar monitoring agent (já vem instalado)
sudo systemctl enable oracle-cloud-agent
sudo systemctl start oracle-cloud-agent

# Ver métricas via CLI
oci monitoring metric-data summarize-metrics-data \
  --compartment-id "xxx" \
  --namespace "oci_computeagent" \
  --query-text "CpuUtilization[1m].mean()"
```

### **Alertas**
```bash
# Criar alarme para CPU alta
oci monitoring alarm create \
  --compartment-id "xxx" \
  --display-name "CPU Alta - retech-core" \
  --metric-compartment-id "xxx" \
  --namespace "oci_computeagent" \
  --query "CpuUtilization[1m].mean() > 80" \
  --severity "CRITICAL" \
  --destinations '["ocid1.onstopic.oc1.xxx"]'
```

---

## 🔄 **CI/CD COM GITHUB ACTIONS:**

### **Workflow GitHub (.github/workflows/deploy-oci.yml)**
```yaml
name: Deploy to Oracle Cloud

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup OCI CLI
        uses: oracle-actions/setup-oci-cli@v1
        with:
          user_ocid: ${{ secrets.OCI_USER_OCID }}
          tenancy_ocid: ${{ secrets.OCI_TENANCY_OCID }}
          fingerprint: ${{ secrets.OCI_FINGERPRINT }}
          api_key: ${{ secrets.OCI_API_KEY }}
          region: sa-saopaulo-1
      
      - name: Deploy via SSH
        run: |
          echo "${{ secrets.SSH_PRIVATE_KEY }}" > key.pem
          chmod 600 key.pem
          
          ssh -i key.pem deploy@${{ secrets.OCI_VM_IP }} << 'EOF'
            cd /app
            git pull origin main
            docker-compose down
            docker-compose up -d --build
            docker-compose logs -f --tail=50
          EOF
```

---

## 💵 **ESTIMATIVA DE CUSTOS:**

### **Always Free (R$ 0,00/mês)**
```
✅ 1 VM ARM (4 cores, 24GB RAM)
✅ 200GB Block Storage
✅ 10TB bandwidth
✅ Load Balancer
✅ Logs + Monitoring

TOTAL: R$ 0,00/mês
```

### **Se precisar expandir (PAGO)**
```
💰 VM adicional (4 cores, 24GB):
   ~R$ 40-60/mês

💰 Block Storage adicional (100GB):
   ~R$ 10/mês

💰 Autonomous Database (se não usar free tier):
   ~R$ 80-120/mês

TOTAL estimado: R$ 130-190/mês (com expansão)
```

### **Comparação com Railway**
```
Railway (atual): $5-10/mês (~R$ 25-50)
Oracle Free:     R$ 0/mês
Oracle Pago:     R$ 130-190/mês (muito mais recursos)
```

---

## 📁 **ESTRUTURA DE SCRIPTS PROPOSTA:**

```bash
scripts/oracle/
├── setup/
│   ├── 00-install-oci-cli.sh     # Instala OCI CLI na máquina local
│   ├── 01-configure-oci.sh       # Configura credentials (interativo)
│   ├── 02-create-vm.sh           # Cria VM + rede (interativo, mostra IP)
│   ├── 03-setup-firewall.sh      # Configura security lists (80, 443, 22)
│   └── 04-install-docker.sh      # SSH na VM e instala Docker
├── deploy/
│   ├── deploy-full.sh            # Deploy completo (services + app)
│   ├── deploy-services.sh        # Apenas Redis/Mongo
│   ├── deploy-app.sh             # Apenas backend
│   └── update-env.sh             # Atualiza variáveis (.env remoto)
├── dns/
│   ├── get-public-ip.sh          # Mostra IP público da VM
│   ├── cloudflare-instructions.sh # Instruções para Cloudflare
│   └── verify-dns.sh             # Verifica propagação DNS
├── monitoring/
│   ├── setup-logs.sh             # Configura OCI Logging
│   ├── view-logs.sh              # Tail logs em tempo real
│   ├── setup-alerts.sh           # Configura alertas (email/Slack)
│   └── dashboard.sh              # Métricas em tempo real (CLI)
├── scale/
│   ├── scale-cpu.sh              # Aumenta/diminui vCPU (interativo)
│   ├── scale-memory.sh           # Aumenta/diminui RAM (interativo)
│   ├── add-storage.sh            # Adiciona block volume
│   └── check-costs.sh            # Verifica custos atuais
├── backup/
│   ├── backup-now.sh             # Backup manual (MongoDB + volumes)
│   ├── setup-auto-backup.sh      # Configura backup automático
│   └── restore.sh                # Restaura de backup (interativo)
├── ci-cd/
│   ├── setup-github-actions.sh   # Configura secrets no GitHub
│   ├── test-deploy.sh            # Testa deploy sem push
│   └── rollback.sh               # Rollback para versão anterior
└── config/
    ├── .env.production           # Variáveis de produção (local)
    ├── docker-compose.oracle.yml # Compose otimizado para Oracle
    ├── nginx.conf                # Reverse proxy
    └── cloud-init.sh             # Script de inicialização da VM
```

---

## 🎯 **FLUXO DE USO (INTERATIVO):**

### **1. Setup Inicial (Uma vez apenas)**
```bash
# Na sua máquina local
cd scripts/oracle/setup

# Passo 1: Instalar OCI CLI
./00-install-oci-cli.sh
# Output: ✅ OCI CLI instalado! Versão: 3.x.x

# Passo 2: Configurar credentials (interativo)
./01-configure-oci.sh
# Perguntas:
# - User OCID? (copie do Oracle Console)
# - Tenancy OCID? (copie do Oracle Console)
# - Region? (digite: sa-saopaulo-1)
# - Gerar chave API? (Y)
# Output: ✅ Configuração salva em ~/.oci/config
#         ✅ Chave API gerada: ~/.oci/oci_api_key.pem
#         📋 Copie a chave pública e adicione no Oracle Console!

# Passo 3: Criar VM (interativo)
./02-create-vm.sh
# Perguntas:
# - Nome da VM? (ex: retech-core-prod)
# - Shape? (1: Free Tier ARM | 2: Micro AMD) → Digite: 1
# - Storage? (padrão: 50GB) → Enter
# - SSH Key? (usa ~/.ssh/id_rsa.pub) → Enter
# Output: 
#   🔄 Criando VM na região São Paulo...
#   ✅ VM criada com sucesso!
#   📋 IP Público: 150.230.45.10 ← GUARDAR ISSO!
#   📋 SSH: ssh ubuntu@150.230.45.10
#   💾 Config salva em: ~/.retech/oracle-vm.json

# Passo 4: Configurar firewall
./03-setup-firewall.sh
# Output:
#   ✅ Porta 80 (HTTP) aberta
#   ✅ Porta 443 (HTTPS) aberta
#   ✅ Porta 22 (SSH) aberta (apenas seu IP)
#   ✅ Porta 8080 (API) aberta

# Passo 5: Instalar Docker na VM
./04-install-docker.sh
# Output:
#   🔄 Conectando via SSH em 150.230.45.10...
#   🔄 Instalando Docker...
#   🔄 Instalando Docker Compose...
#   ✅ Docker instalado! Versão: 24.x
#   ✅ Docker Compose instalado! Versão: 2.x
```

---

### **2. Deploy da Aplicação**
```bash
cd scripts/oracle/deploy

# Deploy completo (primeira vez)
./deploy-full.sh
# Perguntas:
# - Environment? (1: production | 2: staging) → Digite: 1
# - MongoDB Password? (digite uma senha forte)
# - JWT Secret? (auto-gerado ou digite)
# - API Key Secret? (9gJlYXwSR1kfAv8Dh4mHRb/WGJKpLV5v+NYDsNFWTJ8=)
# Output:
#   🔄 Enviando arquivos via SSH...
#   🔄 Criando .env remoto...
#   🔄 Iniciando Redis...
#   ✅ Redis: UP (port 6379)
#   🔄 Iniciando MongoDB...
#   ✅ MongoDB: UP (port 27017)
#   🔄 Building backend...
#   🔄 Iniciando backend...
#   ✅ Backend: UP (port 8080)
#   
#   ✅ DEPLOY COMPLETO!
#   📋 API URL: http://150.230.45.10:8080
#   📋 Health: http://150.230.45.10:8080/health
```

---

### **3. Configurar DNS no Cloudflare**
```bash
cd scripts/oracle/dns

# Obter IP público
./get-public-ip.sh
# Output:
#   📋 IP Público da VM: 150.230.45.10
#   📋 Copie este IP e adicione no Cloudflare!

# Ver instruções
./cloudflare-instructions.sh
# Output:
#   📝 PASSOS PARA CLOUDFLARE:
#   
#   1. Acesse: https://dash.cloudflare.com
#   2. Selecione: theretech.com.br
#   3. Vá em: DNS > Records
#   4. Adicione:
#      Type: A
#      Name: core  (ou api, backend, etc)
#      IPv4: 150.230.45.10
#      Proxy: OFF (🔴) ← IMPORTANTE!
#      TTL: Auto
#   5. Salve!
#   
#   ⏱️ Propagação: 1-5 minutos
#   ✅ URL final: https://core.theretech.com.br

# Verificar DNS
./verify-dns.sh core.theretech.com.br
# Output:
#   🔄 Verificando DNS para core.theretech.com.br...
#   ✅ DNS resolvido: 150.230.45.10
#   ✅ Porta 80: Aberta
#   ✅ Porta 443: Aberta
#   ✅ API respondendo: {"status":"ok","version":"1.0.0"}
```

---

### **4. Atualizar Variáveis de Ambiente**
```bash
cd scripts/oracle/deploy

# Editar variáveis
nano config/.env.production
# Edite:
# MONGO_URI=mongodb://mongo:27017
# REDIS_URL=redis://redis:6379
# JWT_ACCESS_SECRET=xxx
# etc...

# Enviar para VM
./update-env.sh
# Output:
#   🔄 Enviando .env para VM...
#   ✅ Arquivo enviado!
#   🔄 Reiniciando serviços...
#   ✅ Backend reiniciado com novas variáveis!
```

---

### **5. Ver Logs (Tempo Real)**
```bash
cd scripts/oracle/monitoring

# Ver logs do backend
./view-logs.sh backend
# Output (tail -f):
#   [2025-10-27 10:30:45] INFO: Server starting on :8080
#   [2025-10-27 10:30:46] INFO: MongoDB connected
#   [2025-10-27 10:30:46] INFO: Redis connected
#   [2025-10-27 10:30:47] INFO: ✅ Server ready!

# Ver logs do Redis
./view-logs.sh redis

# Ver logs do MongoDB
./view-logs.sh mongo

# Ver TODOS os logs
./view-logs.sh all
```

---

### **6. Escalar Recursos (Interativo)**
```bash
cd scripts/oracle/scale

# Aumentar CPU
./scale-cpu.sh
# Perguntas:
# - CPU atual: 4 cores
# - Nova quantidade? (4-64) → Digite: 8
# - Confirma? (Y/n) → Y
# - Custo: R$ 0/mês (dentro do free tier) OU R$ 40/mês
# Output:
#   ⚠️ Esta ação vai REINICIAR a VM!
#   🔄 Parando VM...
#   🔄 Alterando shape para 8 cores...
#   🔄 Iniciando VM...
#   ✅ VM escalada! Agora tem 8 cores.
#   💰 Custo: R$ 0/mês (free tier)

# Adicionar storage
./add-storage.sh
# Perguntas:
# - Tamanho? (GB) → Digite: 100
# - Custo: R$ 10/mês
# - Confirma? (Y/n) → Y
# Output:
#   🔄 Criando volume de 100GB...
#   ✅ Volume criado: ocid1.volume.xxx
#   🔄 Anexando à VM...
#   ✅ Volume anexado em /dev/sdb
#   📋 Próximo passo: SSH na VM e montar o volume
#   💰 Custo adicional: R$ 10/mês
```

---

### **7. CI/CD Automático**
```bash
cd scripts/oracle/ci-cd

# Setup GitHub Actions (uma vez)
./setup-github-actions.sh
# Perguntas:
# - GitHub Token? (crie em Settings > Developer Settings)
# - Branch para deploy? (padrão: main) → Enter
# - Notificações? (1: Slack | 2: Discord | 3: Email) → Digite: 1
# - Slack Webhook? (cole a URL)
# Output:
#   ✅ Secrets adicionados no GitHub:
#      - OCI_USER_OCID
#      - OCI_TENANCY_OCID
#      - OCI_API_KEY
#      - SSH_PRIVATE_KEY
#      - VM_PUBLIC_IP
#      - SLACK_WEBHOOK
#   ✅ Workflow criado: .github/workflows/deploy-oci.yml
#   
#   📋 Agora, ao fazer push em main:
#   1. GitHub Actions detecta push
#   2. Conecta via SSH na VM Oracle
#   3. Faz git pull
#   4. Rebuild dos containers
#   5. Restart da aplicação
#   6. Health check
#   7. Notifica no Slack (✅ ou ❌)
```

---

## 🚀 **CASOS DE USO AVANÇADOS**

### **Cenário 1: Novo Microservice (ex: retech-payments)**
```bash
cd scripts/oracle/microservices

# Criar novo microservice
./create-microservice.sh
# Perguntas:
# - Nome do microservice? → Digite: retech-payments
# - Porta? → Digite: 8081
# - Banco de dados? (1: Compartilhar MongoDB | 2: Novo PostgreSQL) → Digite: 2
# - Redis? (1: Compartilhar | 2: Nova instância) → Digite: 1
# - Repositório GitHub? → Digite: github.com/theretech/retech-payments
# Output:
#   🔄 Criando estrutura para retech-payments...
#   ✅ Diretório criado: /app/retech-payments
#   ✅ docker-compose.payments.yml criado
#   ✅ .env.payments criado
#   🔄 Adicionando PostgreSQL...
#   ✅ PostgreSQL: UP (port 5432)
#   🔄 Configurando Nginx reverse proxy...
#   ✅ Nginx atualizado:
#      - core.theretech.com.br → :8080 (retech-core)
#      - payments.theretech.com.br → :8081 (retech-payments) ← NOVO!
#   
#   📋 Próximos passos:
#   1. Adicione no Cloudflare:
#      Type: A
#      Name: payments
#      IPv4: 150.230.45.10
#   
#   2. Faça deploy:
#      ./deploy-microservice.sh retech-payments
```

---

### **Cenário 2: Ambiente de Staging**
```bash
cd scripts/oracle/environments

# Criar ambiente staging
./create-staging.sh
# Perguntas:
# - VM separada ou mesma VM? (1: Separada | 2: Mesma - porta diferente) → Digite: 2
# - Porta staging? → Digite: 8090
# - Branch GitHub? → Digite: staging
# Output:
#   ✅ Ambiente staging configurado!
#   📋 URLs:
#      - Produção: https://core.theretech.com.br (porta 8080)
#      - Staging: https://staging.theretech.com.br (porta 8090)
#   
#   📋 CI/CD:
#      - Push em 'main' → Deploy produção
#      - Push em 'staging' → Deploy staging
```

---

### **Cenário 3: Load Balancer (Alta Disponibilidade)**
```bash
cd scripts/oracle/load-balancer

# Setup load balancer (2+ VMs)
./setup-lb.sh
# Perguntas:
# - Quantidade de VMs backend? (padrão: 2) → Digite: 2
# - Shape das VMs? (1: Free ARM | 2: Pago) → Digite: 1
# Output:
#   🔄 Criando VM 1...
#   ✅ VM 1: 150.230.45.10
#   🔄 Criando VM 2...
#   ✅ VM 2: 150.230.45.11
#   🔄 Criando Load Balancer...
#   ✅ Load Balancer: 150.230.45.20 (IP virtual)
#   
#   📋 Configuração Cloudflare:
#      Type: A
#      Name: core
#      IPv4: 150.230.45.20 ← Use o IP do LB, não das VMs!
#   
#   ✅ Traffic:
#      - 50% → VM 1 (150.230.45.10)
#      - 50% → VM 2 (150.230.45.11)
#   
#   💰 Custo: R$ 0/mês (LB grátis no Always Free)
```

---

### **Cenário 4: Rollback Rápido**
```bash
cd scripts/oracle/ci-cd

# Rollback para versão anterior
./rollback.sh
# Perguntas:
# - Versão? (mostra últimas 5) → Digite: 2
# Output:
#   📋 Versões disponíveis:
#   1. v1.2.5 (atual) - 27/10/2025 14:30
#   2. v1.2.4 - 27/10/2025 10:15  ← ESTA
#   3. v1.2.3 - 26/10/2025 18:45
#   4. v1.2.2 - 26/10/2025 12:00
#   5. v1.2.1 - 25/10/2025 20:30
#   
#   ⚠️ Rollback para v1.2.4?
#   🔄 Parando containers...
#   🔄 Checkout git para tag v1.2.4...
#   🔄 Rebuilding...
#   🔄 Reiniciando...
#   ✅ Rollback completo!
#   📋 Versão atual: v1.2.4
```

---

### **Cenário 5: Backup Automático**
```bash
cd scripts/oracle/backup

# Configurar backup automático diário
./setup-auto-backup.sh
# Perguntas:
# - Horário do backup? (ex: 03:00) → Digite: 03:00
# - Retenção? (dias) → Digite: 30
# - Destino? (1: Object Storage Oracle | 2: S3 | 3: Local VM) → Digite: 1
# Output:
#   ✅ Cron job criado: 0 3 * * * /app/scripts/backup.sh
#   ✅ Object Storage configurado (20GB free)
#   📋 Backups:
#      - MongoDB: Dump completo (.gz)
#      - Volumes: Snapshot
#      - Configs: .env, docker-compose.yml
#   📋 Retenção: 30 dias (últimos 30 backups mantidos)
#   💰 Custo: R$ 0/mês (dentro do free tier)
```

---

### **Cenário 6: Monitoramento Avançado**
```bash
cd scripts/oracle/monitoring

# Setup alertas críticos
./setup-alerts.sh
# Perguntas:
# - Email para alertas? → Digite: alan@theretech.com.br
# - Slack webhook? (opcional) → Cole URL ou Enter
# - Discord webhook? (opcional) → Enter
# Output:
#   ✅ Alarmes configurados:
#      - 🔴 CPU > 80% (15min)
#      - 🔴 RAM > 90% (10min)
#      - 🔴 Disco > 85% (30min)
#      - 🔴 API não responde (health check fail)
#      - 🔴 MongoDB down
#      - 🔴 Redis down
#   
#   📧 Notificações via:
#      - Email: alan@theretech.com.br
#      - Slack: #alerts channel
#   
#   📊 Dashboard OCI: https://cloud.oracle.com/monitoring

# Ver métricas em tempo real (CLI)
./dashboard.sh
# Output (atualiza a cada 5s):
#   ┌────────────────────────────────────────┐
#   │  📊 RETECH CORE - METRICS DASHBOARD   │
#   ├────────────────────────────────────────┤
#   │  VM: retech-core-prod (150.230.45.10) │
#   │  Region: São Paulo (sa-saopaulo-1)    │
#   │  Uptime: 15 dias, 8 horas             │
#   ├────────────────────────────────────────┤
#   │  🖥️  CPU: 25% (1/4 cores)              │
#   │  💾 RAM: 8.5GB / 24GB (35%)            │
#   │  💿 Disk: 35GB / 200GB (17%)           │
#   │  🌐 Network: ↓ 150 Mbps ↑ 80 Mbps     │
#   ├────────────────────────────────────────┤
#   │  🐳 Docker Containers:                 │
#   │     ✅ retech-api (UP - 15d)           │
#   │     ✅ mongo (UP - 15d)                │
#   │     ✅ redis (UP - 15d)                │
#   │     ✅ nginx (UP - 15d)                │
#   ├────────────────────────────────────────┤
#   │  📊 API Stats (últimas 24h):           │
#   │     Requests: 15.430                   │
#   │     Errors: 23 (0.15%)                 │
#   │     Avg Response: 12ms                 │
#   │     Cache Hit Rate: 92%                │
#   └────────────────────────────────────────┘
#   
#   [Atualiza a cada 5s - Ctrl+C para sair]
```

---

### **Cenário 7: SSL/HTTPS Automático (Let's Encrypt)**
```bash
cd scripts/oracle/ssl

# Setup SSL com Let's Encrypt
./setup-ssl.sh
# Perguntas:
# - Domínio? → Digite: core.theretech.com.br
# - Email? → Digite: alan@theretech.com.br
# Output:
#   ✅ Certbot instalado
#   🔄 Solicitando certificado SSL...
#   ✅ Certificado emitido!
#   🔄 Configurando Nginx...
#   ✅ Nginx configurado para HTTPS
#   🔄 Auto-renovação configurada (cron)
#   
#   📋 URLs:
#      - HTTP: http://core.theretech.com.br (redireciona para HTTPS)
#      - HTTPS: https://core.theretech.com.br ✅
#   
#   📅 Renovação automática: A cada 60 dias
#   💰 Custo: R$ 0/mês (Let's Encrypt gratuito)
```

---

### **Cenário 8: Multi-Região (Redundância Geográfica)**
```bash
cd scripts/oracle/multi-region

# Criar instância em outra região (ex: Chile)
./create-secondary-region.sh
# Perguntas:
# - Região secundária? (1: Chile | 2: EUA | 3: Europa) → Digite: 1
# - Sync em tempo real? (Y/n) → Y
# Output:
#   🔄 Criando VM em Santiago (Chile)...
#   ✅ VM criada: 150.240.50.10
#   🔄 Configurando replicação MongoDB...
#   ✅ MongoDB Replica Set: BR (master) ↔ Chile (slave)
#   🔄 Configurando Redis sync...
#   ✅ Redis replication configurada
#   
#   📋 Latência esperada:
#      - Brasil → Chile: ~40ms
#      - Failover automático se BR cair
#   
#   💰 Custo: R$ 0/mês (2 VMs no free tier)
```

---

### **Cenário 9: Deploy de Hotfix (Urgente)**
```bash
cd scripts/oracle/deploy

# Deploy rápido (pula build, usa imagem Docker Hub)
./hotfix-deploy.sh
# Perguntas:
# - Tag da imagem? (ex: v1.2.6-hotfix) → Digite: v1.2.6-hotfix
# - Fazer backup antes? (Y/n) → Y
# Output:
#   🔄 Criando backup...
#   ✅ Backup salvo: /backups/pre-hotfix-20251027.tar.gz
#   🔄 Pulling imagem: theretech/retech-core:v1.2.6-hotfix
#   🔄 Parando containers...
#   🔄 Iniciando nova versão...
#   🔄 Health check...
#   ✅ Health OK! API respondendo
#   
#   ⏱️ Total: 45 segundos (vs 5 minutos build completo)
#   
#   📋 Rollback rápido (se necessário):
#      ./rollback.sh --version v1.2.5
```

---

### **Cenário 10: Secrets Management (OCI Vault)**
```bash
cd scripts/oracle/secrets

# Armazenar secrets de forma segura
./store-secrets.sh
# Perguntas:
# - Secret name? → Digite: APIKEY_HASH_SECRET
# - Secret value? → Digite: 9gJlYXwSR1kfAv8Dh4mHRb/WGJKpLV5v+NYDsNFWTJ8=
# - Compartment? (usa padrão) → Enter
# Output:
#   ✅ Secret armazenado no OCI Vault
#   📋 Secret ID: ocid1.vaultsecret.oc1.xxx
#   
#   🔄 Atualizando .env remoto para usar vault...
#   ✅ App agora busca secrets do Vault (mais seguro)
#   
#   💰 Custo: R$ 0,03/secret/mês (~R$ 0,30/mês para 10 secrets)

# Rotacionar secret
./rotate-secret.sh APIKEY_HASH_SECRET
# Output:
#   🔄 Gerando novo secret...
#   ✅ Novo secret: xYz123...
#   🔄 Atualizando Vault...
#   🔄 Reiniciando aplicação...
#   ✅ Secret rotacionado com zero downtime!
```

---

### **Cenário 11: Database Migration**
```bash
cd scripts/oracle/database

# Executar migration no MongoDB remoto
./run-migration.sh
# Perguntas:
# - Migration file? → Digite: migrations/001_add_indexes.js
# - Ambiente? (1: production | 2: staging) → Digite: 2
# - Confirma? (Y/n) → Y
# Output:
#   🔄 Conectando em MongoDB staging...
#   🔄 Executando migration: 001_add_indexes.js
#   ✅ Index criado: api_keys.keyId
#   ✅ Index criado: tenants.email
#   ✅ Migration completa!
#   📊 Tempo: 2.3 segundos

# Backup antes de migration perigosa
./migrate-with-backup.sh
# Output:
#   🔄 Criando backup completo...
#   ✅ Backup: /backups/pre-migration-20251027.gz
#   🔄 Executando migration...
#   ✅ Migration OK!
#   📋 Rollback disponível: ./restore-backup.sh
```

---

### **Cenário 12: Custom Domain + SSL**
```bash
cd scripts/oracle/domains

# Adicionar novo domínio
./add-domain.sh
# Perguntas:
# - Domínio? → Digite: api.theretech.com.br
# - Serviço? (1: retech-core | 2: retech-payments) → Digite: 1
# - SSL? (Y/n) → Y
# Output:
#   🔄 Verificando DNS...
#   ⚠️ DNS ainda não resolvido para api.theretech.com.br
#   
#   📋 AÇÃO NECESSÁRIA:
#   Adicione no Cloudflare:
#      Type: A
#      Name: api
#      IPv4: 150.230.45.10
#   
#   Aguardando propagação DNS... (pressione Enter quando pronto)
#   [você adiciona no Cloudflare e pressiona Enter]
#   
#   ✅ DNS resolvido!
#   🔄 Solicitando certificado SSL...
#   ✅ Certificado emitido para api.theretech.com.br
#   🔄 Configurando Nginx...
#   ✅ HTTPS ativo!
#   
#   📋 URLs finais:
#      - https://api.theretech.com.br ✅
#      - https://core.theretech.com.br ✅
```

---

### **Cenário 13: Health Checks & Auto-Recovery**
```bash
cd scripts/oracle/health

# Setup auto-recovery
./setup-auto-recovery.sh
# Output:
#   ✅ Cron job criado: */5 * * * * (a cada 5 min)
#   📋 Verifica:
#      - API respondendo? (curl /health)
#      - MongoDB conectado?
#      - Redis conectado?
#   
#   📋 Se falhar:
#      - Tenta restart do container
#      - Se falhar 3x, notifica admin (email)
#      - Se falhar 5x, reinicia VM completa
#   
#   📊 Logs salvos em: /var/log/health-checks.log
```

---

### **Cenário 14: Resource Monitoring & Auto-Scale**
```bash
cd scripts/oracle/auto-scale

# Setup auto-scaling (se CPU >80% por 10min, aumenta)
./setup-auto-scale.sh
# Perguntas:
# - Habilitar auto-scale? (Y/n) → Y
# - CPU threshold? (%) → Digite: 80
# - Duração? (minutos) → Digite: 10
# - Máximo de cores? → Digite: 8
# Output:
#   ✅ Alarme OCI criado
#   ✅ Webhook configurado
#   📋 Regra:
#      SE CPU > 80% por 10min ENTÃO
#         Aumentar de 4 cores → 8 cores
#      SE CPU < 40% por 30min ENTÃO
#         Diminuir de 8 cores → 4 cores
#   
#   💰 Custo:
#      - 4 cores: R$ 0/mês (free tier)
#      - 8 cores: R$ 40/mês (quando escalar)
#      - Auto volta para 4 cores quando normalizar
```

---

### **Cenário 15: Blue-Green Deployment**
```bash
cd scripts/oracle/deployment-strategies

# Deploy blue-green (zero downtime)
./deploy-blue-green.sh
# Output:
#   📊 Estado atual:
#      - BLUE (produção): porta 8080, versão v1.2.5
#      - GREEN (staging): porta 8081, versão v1.2.6
#   
#   🔄 Testando GREEN...
#   ✅ Health check OK
#   ✅ Smoke tests passaram
#   
#   ⚠️ Trocar tráfego para GREEN? (Y/n) → Y
#   
#   🔄 Atualizando Nginx...
#   🔄 Redirecionando 100% do tráfego para GREEN (8081)...
#   ✅ Tráfego migrado!
#   
#   📋 BLUE ainda rodando (rollback instantâneo se necessário)
#   📋 Aguardando 5 minutos para confirmar estabilidade...
#   
#   [5 minutos depois]
#   ✅ GREEN estável!
#   🔄 Parando BLUE...
#   ✅ Deploy completo com zero downtime!
```

---

### **Cenário 16: Disaster Recovery**
```bash
cd scripts/oracle/disaster-recovery

# Recuperar de desastre total (VM perdida)
./recover-from-disaster.sh
# Perguntas:
# - Backup para restaurar? (lista últimos 10) → Digite: 3
# Output:
#   📋 Backups disponíveis:
#   1. 27/10/2025 03:00 (mais recente)
#   2. 26/10/2025 03:00
#   3. 25/10/2025 03:00  ← ESTE
#   
#   ⚠️ Este processo vai:
#   1. Criar nova VM
#   2. Restaurar MongoDB do backup
#   3. Restaurar volumes
#   4. Reconfigurar tudo
#   
#   Continuar? (Y/n) → Y
#   
#   🔄 Criando nova VM...
#   ✅ Nova VM: 150.230.45.50
#   🔄 Instalando Docker...
#   🔄 Restaurando MongoDB (3.2GB)...
#   🔄 Restaurando volumes (15GB)...
#   🔄 Iniciando serviços...
#   ✅ Tudo restaurado!
#   
#   📋 AÇÃO NECESSÁRIA:
#   Atualize no Cloudflare:
#      core.theretech.com.br
#      IPv4: 150.230.45.10 → 150.230.45.50
#   
#   ⏱️ Downtime total: ~15 minutos
#   ✅ RTO (Recovery Time Objective): <30 min
#   ✅ RPO (Recovery Point Objective): <24h (backup diário)
```

---

### **Cenário 17: Performance Testing**
```bash
cd scripts/oracle/testing

# Teste de carga
./load-test.sh
# Perguntas:
# - Endpoint? → Digite: /cep/01310100
# - Requests/segundo? → Digite: 100
# - Duração? (segundos) → Digite: 60
# Output:
#   🔄 Executando teste de carga...
#   📊 Target: 100 req/s por 60s (6.000 requests total)
#   
#   [Barra de progresso]
#   ██████████████████████████████ 100% (60s)
#   
#   📊 RESULTADOS:
#   ├─ Total requests: 6.000
#   ├─ Successful: 5.997 (99.95%)
#   ├─ Failed: 3 (0.05%)
#   ├─ Avg response: 8ms
#   ├─ P50: 5ms
#   ├─ P95: 15ms
#   ├─ P99: 45ms
#   ├─ Max: 120ms
#   └─ Cache hit rate: 94%
#   
#   📊 Resource Usage:
#   ├─ CPU peak: 65%
#   ├─ RAM peak: 12GB (50%)
#   └─ Network: 80 Mbps
#   
#   ✅ Sistema aguenta 100 req/s facilmente!
#   💡 Recomendação: Pode escalar até 500 req/s na config atual
```

---

## 📁 **ESTRUTURA COMPLETA DE SCRIPTS:**

```bash
scripts/oracle/
├── setup/                        # Setup inicial (uma vez)
│   ├── 00-install-oci-cli.sh
│   ├── 01-configure-oci.sh
│   ├── 02-create-vm.sh
│   ├── 03-setup-firewall.sh
│   └── 04-install-docker.sh
├── deploy/                       # Deploy e updates
│   ├── deploy-full.sh
│   ├── deploy-app.sh
│   ├── deploy-services.sh
│   ├── hotfix-deploy.sh         # Deploy rápido
│   └── update-env.sh
├── microservices/                # Gerenciar múltiplos services
│   ├── create-microservice.sh
│   ├── deploy-microservice.sh
│   └── list-services.sh
├── environments/                 # Staging, QA, etc
│   ├── create-staging.sh
│   ├── create-qa.sh
│   └── switch-env.sh
├── dns/                          # DNS e domínios
│   ├── get-public-ip.sh
│   ├── cloudflare-instructions.sh
│   └── verify-dns.sh
├── ssl/                          # HTTPS e certificados
│   ├── setup-ssl.sh
│   ├── renew-ssl.sh
│   └── add-domain.sh
├── load-balancer/               # Alta disponibilidade
│   ├── setup-lb.sh
│   ├── add-backend.sh
│   └── remove-backend.sh
├── monitoring/                   # Logs e métricas
│   ├── setup-logs.sh
│   ├── view-logs.sh
│   ├── setup-alerts.sh
│   ├── dashboard.sh             # CLI dashboard
│   └── export-metrics.sh
├── scale/                        # Escalabilidade
│   ├── scale-cpu.sh
│   ├── scale-memory.sh
│   ├── add-storage.sh
│   ├── setup-auto-scale.sh      # Auto-scaling
│   └── check-costs.sh
├── backup/                       # Backup e restore
│   ├── backup-now.sh
│   ├── setup-auto-backup.sh
│   ├── restore.sh
│   └── list-backups.sh
├── database/                     # Database operations
│   ├── run-migration.sh
│   ├── migrate-with-backup.sh
│   ├── export-db.sh
│   └── import-db.sh
├── ci-cd/                        # Integração contínua
│   ├── setup-github-actions.sh
│   ├── test-deploy.sh
│   ├── rollback.sh
│   └── deploy-blue-green.sh
├── disaster-recovery/            # Recuperação de desastres
│   ├── recover-from-disaster.sh
│   ├── create-dr-plan.sh
│   └── test-dr.sh
├── multi-region/                 # Redundância geográfica
│   ├── create-secondary-region.sh
│   ├── setup-replication.sh
│   └── failover-test.sh
├── testing/                      # Testes e validação
│   ├── load-test.sh
│   ├── smoke-test.sh
│   └── integration-test.sh
├── secrets/                      # Gestão de secrets
│   ├── store-secret.sh
│   ├── rotate-secret.sh
│   ├── list-secrets.sh
│   └── delete-secret.sh
└── utils/                        # Utilitários
    ├── ssh-connect.sh
    ├── port-forward.sh          # Túnel SSH para MongoDB/Redis
    ├── check-health.sh
    └── estimate-costs.sh
```

---

## 🔧 **SCRIPTS ADICIONAIS ESSENCIAIS:**

### **Port Forwarding (Acesso Local a Serviços Remotos)**
```bash
cd scripts/oracle/utils

# Criar túnel SSH para MongoDB
./port-forward.sh mongo
# Output:
#   🔄 Criando túnel SSH...
#   ✅ Túnel ativo!
#   📋 Conexão local:
#      mongodb://localhost:27017
#   
#   💡 Use MongoDB Compass:
#      - Host: localhost
#      - Port: 27017
#      - Auth: (mesmas credenciais remotas)
#   
#   ⚠️ Túnel ativo enquanto este terminal estiver aberto
#   Ctrl+C para fechar

# Túnel para Redis
./port-forward.sh redis
# Output:
#   ✅ Túnel ativo!
#   📋 Conexão local: redis://localhost:6379
#   💡 Use RedisInsight ou redis-cli
```

---

### **Cost Estimation (Antes de Escalar)**
```bash
cd scripts/oracle/utils

# Estimar custos antes de aumentar recursos
./estimate-costs.sh
# Perguntas:
# - Scenario? (1: Current | 2: Upgrade Plan) → Digite: 2
# - vCPUs? → Digite: 8
# - RAM (GB)? → Digite: 32
# - Storage (GB)? → Digite: 300
# Output:
#   📊 ESTIMATIVA DE CUSTOS:
#   
#   ┌─────────────────────────────────────┐
#   │ Configuração Atual (Free Tier):     │
#   ├─────────────────────────────────────┤
#   │ vCPU: 4 cores ARM                   │
#   │ RAM: 24GB                           │
#   │ Storage: 200GB                      │
#   │ 💰 Custo: R$ 0/mês                  │
#   └─────────────────────────────────────┘
#   
#   ┌─────────────────────────────────────┐
#   │ Configuração Planejada:             │
#   ├─────────────────────────────────────┤
#   │ vCPU: 8 cores ARM (+4)              │
#   │ RAM: 32GB (+8GB)                    │
#   │ Storage: 300GB (+100GB)             │
#   │ 💰 Custo: R$ 65/mês                 │
#   │    ├─ vCPU extra: R$ 40/mês         │
#   │    ├─ RAM extra: R$ 15/mês          │
#   │    └─ Storage: R$ 10/mês            │
#   └─────────────────────────────────────┘
#   
#   📊 Comparação Railway:
#      - Railway equivalente: ~$50/mês (~R$ 250/mês)
#      - Oracle: R$ 65/mês
#      - Economia: R$ 185/mês (74% mais barato)
#      - Benefício: Servidor no Brasil (10x mais rápido)
```

---

## 🎯 **SCRIPTS QUE FALTAM (A DESENVOLVER):**

### **Essenciais**
- [ ] `00-install-oci-cli.sh` - Instalar OCI CLI
- [ ] `01-configure-oci.sh` - Setup credentials
- [ ] `02-create-vm.sh` - Criar VM interativo
- [ ] `deploy-full.sh` - Deploy completo
- [ ] `get-public-ip.sh` - Mostrar IP para Cloudflare
- [ ] `view-logs.sh` - Tail logs remotos
- [ ] `rollback.sh` - Voltar versão

### **Importantes**
- [ ] `create-microservice.sh` - Novo serviço
- [ ] `setup-ssl.sh` - HTTPS automático
- [ ] `setup-auto-backup.sh` - Backup diário
- [ ] `setup-alerts.sh` - Notificações
- [ ] `dashboard.sh` - Métricas CLI

### **Avançados**
- [ ] `setup-lb.sh` - Load balancer
- [ ] `deploy-blue-green.sh` - Zero downtime
- [ ] `setup-auto-scale.sh` - Auto-scaling
- [ ] `create-secondary-region.sh` - Multi-região
- [ ] `store-secret.sh` - OCI Vault

---

## 💰 **ESTIMATIVA FINAL DE CUSTOS:**

### **Cenário 1: Free Tier Completo (Recomendado)**
```
✅ 1 VM ARM (4 cores, 24GB RAM)
✅ 200GB Block Storage
✅ MongoDB + Redis + Backend (Docker)
✅ Load Balancer (1 instância)
✅ SSL gratuito (Let's Encrypt)
✅ Backup automático (Object Storage 20GB)
✅ Monitoramento + Logs

💰 TOTAL: R$ 0,00/mês
```

### **Cenário 2: Expansão Moderada**
```
✅ 1 VM ARM (8 cores, 32GB RAM) ← Upgrade
✅ 300GB Block Storage ← +100GB
✅ Tudo do Cenário 1

💰 TOTAL: R$ 65/mês
   ├─ vCPU: R$ 40/mês
   ├─ RAM: R$ 15/mês
   └─ Storage: R$ 10/mês
```

### **Cenário 3: Alta Disponibilidade (Multi-Região)**
```
✅ 2 VMs ARM (4 cores, 24GB cada) - BR + Chile
✅ MongoDB Replica Set (2 regiões)
✅ Redis replication
✅ Load Balancer global
✅ 400GB Storage total

💰 TOTAL: R$ 0/mês (2 VMs no free tier)
```

---

**TODOS ESSES CENÁRIOS PODEM SER AUTOMATIZADOS VIA SCRIPTS .sh!** 🚀

---

| Funcionalidade | Railway (GUI) | Oracle (CLI/Script) |
|----------------|---------------|---------------------|
| Deploy automático | ✅ GitHub integration | ✅ GitHub Actions |
| Variáveis de ambiente | ✅ Interface web | ✅ .env + OCI Vault |
| Logs | ✅ Web dashboard | ✅ OCI Logging + CLI |
| Métricas | ✅ Gráficos web | ✅ OCI Monitoring + CLI |
| Escalabilidade | ✅ Slider web | ✅ Scripts bash |
| Rollback | ✅ 1 click | ✅ Script + Git tag |
| Secrets | ✅ Interface web | ✅ OCI Vault |
| Custom domains | ✅ Interface web | ✅ DNS + Load Balancer |

---

## 🎯 **PRÓXIMOS PASSOS:**

1. ✅ Criar conta Oracle Cloud
2. ✅ Configurar OCI CLI localmente
3. ✅ Desenvolver scripts de automação
4. ✅ Testar deploy em ambiente de staging
5. ✅ Migrar produção gradualmente
6. ✅ Documentar processo completo

---

## 📚 **DOCUMENTAÇÃO OFICIAL:**

- OCI CLI: https://docs.oracle.com/en-us/iaas/tools/oci-cli/latest/
- Compute: https://docs.oracle.com/en-us/iaas/Content/Compute/home.htm
- Always Free: https://www.oracle.com/cloud/free/
- GitHub Actions: https://github.com/oracle-actions/

---

**⚠️ IMPORTANTE:**
- Always Free Tier é realmente GRATUITO para sempre
- Sem cartão de crédito necessário (mas recomendado para upgrade futuro)
- Região São Paulo disponível para Always Free
- Recursos podem ser expandidos pagando sob demanda

---

**Status:** 📝 Documentação de pesquisa - Aguardando aprovação para desenvolver scripts

