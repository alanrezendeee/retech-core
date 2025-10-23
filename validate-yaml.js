const fs = require('fs');
const http = require('http');

// Baixar YAML da API
http.get('http://localhost:8080/openapi.yaml', (res) => {
  let data = '';
  
  res.on('data', (chunk) => {
    data += chunk;
  });
  
  res.on('end', () => {
    console.log('📄 YAML recebido, analisando...\n');
    
    // Dividir em linhas
    const lines = data.split('\n');
    
    // Verificar linhas suspeitas
    lines.forEach((line, index) => {
      const lineNum = index + 1;
      
      // Verificar indentação inconsistente
      if (line.match(/^\s+\w+:/)) {
        const spaces = line.match(/^(\s+)/);
        if (spaces && spaces[1].length % 2 !== 0) {
          console.log(`⚠️  Linha ${lineNum}: Indentação ímpar (${spaces[1].length} espaços)`);
          console.log(`    ${line}`);
        }
      }
      
      // Linha 159 específica
      if (lineNum === 159) {
        console.log(`\n🔍 Linha 159 (a problemática):`);
        console.log(`    Conteúdo: "${line}"`);
        console.log(`    Tamanho: ${line.length} chars`);
        console.log(`    Espaços no início: ${line.match(/^(\s*)/)[1].length}`);
        
        // Mostrar contexto
        console.log(`\n📋 Contexto (linhas 155-165):`);
        for (let i = 154; i < 165; i++) {
          const marker = i === 158 ? ' ⚠️ ' : '    ';
          console.log(`${marker}${i + 1}: ${lines[i]}`);
        }
      }
    });
    
    console.log('\n✅ Análise concluída!');
  });
}).on('error', (err) => {
  console.error('❌ Erro:', err.message);
});

