#!/usr/bin/env node

/**
 * 📊 Gerador de Gráfico de Resultados de Rate Limiting
 * Lê rate-limit-test-results.json e gera visualização ASCII
 */

const fs = require('fs');
const path = require('path');

// Cores ANSI
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
};

// Carregar resultados
const resultsFile = path.join(__dirname, '../rate-limit-test-results.json');

if (!fs.existsSync(resultsFile)) {
  console.error(`${colors.red}❌ Arquivo de resultados não encontrado: ${resultsFile}${colors.reset}`);
  console.error(`   Execute primeiro: ./scripts/test-rate-limit.sh`);
  process.exit(1);
}

const results = JSON.parse(fs.readFileSync(resultsFile, 'utf8'));

console.log('');
console.log(`${colors.bright}${colors.cyan}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${colors.reset}`);
console.log(`${colors.bright}${colors.cyan}   📊 RELATÓRIO DE TESTES DE RATE LIMITING${colors.reset}`);
console.log(`${colors.bright}${colors.cyan}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${colors.reset}`);
console.log('');

// Estatísticas gerais
const totalTests = results.length;
const passedTests = results.filter(r => r.passed === true || r.passed === 'true').length;
const failedTests = totalTests - passedTests;

console.log(`${colors.bright}RESUMO GERAL:${colors.reset}`);
console.log(`  Total de testes:    ${totalTests}`);
console.log(`  ✅ Passaram:        ${colors.green}${passedTests}${colors.reset}`);
console.log(`  ❌ Falharam:        ${colors.red}${failedTests}${colors.reset}`);
console.log(`  Taxa de sucesso:   ${colors.cyan}${((passedTests / totalTests) * 100).toFixed(1)}%${colors.reset}`);
console.log('');

// Gráfico de barras para cada cenário
console.log(`${colors.bright}DETALHES POR CENÁRIO:${colors.reset}`);
console.log('');

results.forEach((result, index) => {
  const isPassed = result.passed === true || result.passed === 'true';
  const statusIcon = isPassed ? '✅' : '❌';
  const statusColor = isPassed ? colors.green : colors.red;
  
  console.log(`${colors.bright}${index + 1}. ${result.scenario}${colors.reset} ${statusIcon}`);
  console.log(`   Limite esperado: ${result.expectedLimit} requests`);
  console.log(`   Requests feitas: ${result.requestsMade}`);
  console.log('');
  
  // Gráfico de barras
  const maxWidth = 50;
  const successBar = '█'.repeat(Math.floor((result.success / result.requestsMade) * maxWidth));
  const rateLimitedBar = '█'.repeat(Math.floor((result.rateLimited / result.requestsMade) * maxWidth));
  const errorBar = '█'.repeat(Math.floor((result.errors / result.requestsMade) * maxWidth));
  
  console.log(`   ${colors.green}✅ Sucesso (${result.success}):${colors.reset}`);
  console.log(`      ${colors.green}${successBar}${colors.reset} ${((result.success / result.requestsMade) * 100).toFixed(1)}%`);
  console.log('');
  
  console.log(`   ${colors.red}🚫 Rate Limited (${result.rateLimited}):${colors.reset}`);
  console.log(`      ${colors.red}${rateLimitedBar}${colors.reset} ${((result.rateLimited / result.requestsMade) * 100).toFixed(1)}%`);
  console.log('');
  
  if (result.errors > 0) {
    console.log(`   ${colors.yellow}⚠️  Erros (${result.errors}):${colors.reset}`);
    console.log(`      ${colors.yellow}${errorBar}${colors.reset} ${((result.errors / result.requestsMade) * 100).toFixed(1)}%`);
    console.log('');
  }
  
  if (result.first429) {
    console.log(`   🎯 Primeiro 429 na request: #${result.first429}`);
  } else {
    console.log(`   ${colors.yellow}⚠️  Nenhum 429 recebido${colors.reset}`);
  }
  
  // Análise
  console.log('');
  console.log(`   ${colors.bright}ANÁLISE:${colors.reset}`);
  
  if (isPassed) {
    console.log(`   ${colors.green}✅ Rate limit funcionou corretamente!${colors.reset}`);
    console.log(`      - ${result.success} requests permitidas (≤ ${result.expectedLimit})`);
    console.log(`      - ${result.rateLimited} requests bloqueadas`);
    if (result.first429 === result.expectedLimit + 1) {
      console.log(`      - Bloqueio ocorreu exatamente após o limite`);
    }
  } else {
    console.log(`   ${colors.red}❌ Problema detectado:${colors.reset}`);
    
    if (result.success > result.expectedLimit) {
      console.log(`      - Permitiu MAIS requests que o limite (${result.success} > ${result.expectedLimit})`);
    }
    
    if (result.rateLimited === 0) {
      console.log(`      - Nenhuma request foi bloqueada (429 não retornado)`);
    }
    
    if (result.first429 && result.first429 !== result.expectedLimit + 1) {
      console.log(`      - Bloqueio ocorreu na request #${result.first429} (esperado: #${result.expectedLimit + 1})`);
    }
  }
  
  console.log('');
  console.log(`${colors.cyan}${'─'.repeat(65)}${colors.reset}`);
  console.log('');
});

// Gráfico de pizza ASCII (aproximado)
console.log(`${colors.bright}DISTRIBUIÇÃO GERAL:${colors.reset}`);
console.log('');

const totalRequests = results.reduce((sum, r) => sum + r.requestsMade, 0);
const totalSuccess = results.reduce((sum, r) => sum + r.success, 0);
const totalRateLimited = results.reduce((sum, r) => sum + r.rateLimited, 0);
const totalErrors = results.reduce((sum, r) => sum + r.errors, 0);

const successPercent = ((totalSuccess / totalRequests) * 100).toFixed(1);
const rateLimitedPercent = ((totalRateLimited / totalRequests) * 100).toFixed(1);
const errorsPercent = ((totalErrors / totalRequests) * 100).toFixed(1);

console.log(`  Total de requests: ${totalRequests}`);
console.log('');
console.log(`  ${colors.green}✅ Sucesso:       ${totalSuccess.toString().padEnd(5)} (${successPercent}%)${colors.reset}`);
console.log(`  ${colors.red}🚫 Rate Limited: ${totalRateLimited.toString().padEnd(5)} (${rateLimitedPercent}%)${colors.reset}`);
console.log(`  ${colors.yellow}⚠️  Erros:        ${totalErrors.toString().padEnd(5)} (${errorsPercent}%)${colors.reset}`);
console.log('');

// Gráfico de barras horizontal
const barWidth = 40;
const successBarWidth = Math.floor((totalSuccess / totalRequests) * barWidth);
const rateLimitedBarWidth = Math.floor((totalRateLimited / totalRequests) * barWidth);
const errorsBarWidth = Math.floor((totalErrors / totalRequests) * barWidth);

console.log(`  [${colors.green}${'█'.repeat(successBarWidth)}${colors.red}${'█'.repeat(rateLimitedBarWidth)}${colors.yellow}${'█'.repeat(errorsBarWidth)}${colors.reset}${'░'.repeat(barWidth - successBarWidth - rateLimitedBarWidth - errorsBarWidth)}]`);
console.log('');

// Conclusão
console.log(`${colors.bright}${colors.cyan}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${colors.reset}`);
console.log(`${colors.bright}CONCLUSÃO:${colors.reset}`);
console.log('');

if (failedTests === 0) {
  console.log(`${colors.green}✅ TODOS OS TESTES PASSARAM!${colors.reset}`);
  console.log(`   O sistema de rate limiting está funcionando corretamente.`);
} else {
  console.log(`${colors.red}❌ ALGUNS TESTES FALHARAM!${colors.reset}`);
  console.log(`   ${failedTests} de ${totalTests} cenários apresentaram problemas.`);
  console.log('');
  console.log(`${colors.yellow}   AÇÕES RECOMENDADAS:${colors.reset}`);
  console.log(`   1. Verificar o código em internal/middleware/rate_limiter.go`);
  console.log(`   2. Conferir se o middleware está aplicado corretamente`);
  console.log(`   3. Verificar se os limites estão sendo lidos do banco de dados`);
  console.log(`   4. Conferir logs do backend para erros`);
}

console.log('');
console.log(`${colors.bright}${colors.cyan}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${colors.reset}`);
console.log('');

// Exportar relatório HTML (opcional)
const htmlReport = generateHTMLReport(results);
const htmlFile = path.join(__dirname, '../rate-limit-test-report.html');
fs.writeFileSync(htmlFile, htmlReport);

console.log(`📄 Relatório HTML gerado: ${htmlFile}`);
console.log(`   Abra em: file://${htmlFile}`);
console.log('');

// Função para gerar relatório HTML
function generateHTMLReport(results) {
  const totalRequests = results.reduce((sum, r) => sum + r.requestsMade, 0);
  const totalSuccess = results.reduce((sum, r) => sum + r.success, 0);
  const totalRateLimited = results.reduce((sum, r) => sum + r.rateLimited, 0);
  
  return `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Relatório de Testes - Rate Limiting</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 40px 20px;
            min-height: 100vh;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }
        h1 {
            color: #667eea;
            margin-bottom: 10px;
            font-size: 32px;
        }
        .subtitle {
            color: #888;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 24px;
            border-radius: 12px;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
        }
        .stat-card h3 {
            font-size: 14px;
            opacity: 0.9;
            margin-bottom: 8px;
        }
        .stat-card .value {
            font-size: 36px;
            font-weight: bold;
        }
        .charts {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
            gap: 30px;
            margin-bottom: 40px;
        }
        .chart-container {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 12px;
        }
        .scenario {
            background: #f8f9fa;
            padding: 24px;
            border-radius: 12px;
            margin-bottom: 20px;
        }
        .scenario h3 {
            color: #333;
            margin-bottom: 16px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .badge {
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
        }
        .badge.passed { background: #10b981; color: white; }
        .badge.failed { background: #ef4444; color: white; }
        .metric {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #e5e7eb;
        }
        .metric:last-child { border-bottom: none; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Relatório de Testes - Rate Limiting</h1>
        <p class="subtitle">Retech Core API - ${new Date().toLocaleString('pt-BR')}</p>
        
        <div class="summary">
            <div class="stat-card">
                <h3>Total de Testes</h3>
                <div class="value">${results.length}</div>
            </div>
            <div class="stat-card">
                <h3>Requests Totais</h3>
                <div class="value">${totalRequests}</div>
            </div>
            <div class="stat-card">
                <h3>Taxa de Sucesso</h3>
                <div class="value">${((results.filter(r => r.passed).length / results.length) * 100).toFixed(0)}%</div>
            </div>
        </div>
        
        <div class="charts">
            <div class="chart-container">
                <h3>Distribuição de Respostas</h3>
                <canvas id="pieChart"></canvas>
            </div>
            <div class="chart-container">
                <h3>Resultados por Cenário</h3>
                <canvas id="barChart"></canvas>
            </div>
        </div>
        
        <h2 style="margin-bottom: 20px;">Detalhes por Cenário</h2>
        
        ${results.map((r, i) => `
            <div class="scenario">
                <h3>
                    ${i + 1}. ${r.scenario}
                    <span class="badge ${r.passed ? 'passed' : 'failed'}">
                        ${r.passed ? '✅ PASSOU' : '❌ FALHOU'}
                    </span>
                </h3>
                <div class="metric">
                    <span>Limite Esperado</span>
                    <strong>${r.expectedLimit} req</strong>
                </div>
                <div class="metric">
                    <span>Requests Feitas</span>
                    <strong>${r.requestsMade}</strong>
                </div>
                <div class="metric">
                    <span>✅ Sucesso</span>
                    <strong style="color: #10b981">${r.success}</strong>
                </div>
                <div class="metric">
                    <span>🚫 Rate Limited</span>
                    <strong style="color: #ef4444">${r.rateLimited}</strong>
                </div>
                <div class="metric">
                    <span>⚠️ Erros</span>
                    <strong style="color: #f59e0b">${r.errors}</strong>
                </div>
                ${r.first429 ? `
                <div class="metric">
                    <span>🎯 Primeiro 429</span>
                    <strong>Request #${r.first429}</strong>
                </div>
                ` : ''}
            </div>
        `).join('')}
    </div>
    
    <script>
        // Pie Chart
        new Chart(document.getElementById('pieChart'), {
            type: 'pie',
            data: {
                labels: ['✅ Sucesso', '🚫 Rate Limited', '⚠️ Erros'],
                datasets: [{
                    data: [${totalSuccess}, ${totalRateLimited}, ${results.reduce((sum, r) => sum + r.errors, 0)}],
                    backgroundColor: ['#10b981', '#ef4444', '#f59e0b']
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: { position: 'bottom' }
                }
            }
        });
        
        // Bar Chart
        new Chart(document.getElementById('barChart'), {
            type: 'bar',
            data: {
                labels: ${JSON.stringify(results.map(r => r.scenario.substring(0, 20)))},
                datasets: [
                    {
                        label: 'Sucesso',
                        data: ${JSON.stringify(results.map(r => r.success))},
                        backgroundColor: '#10b981'
                    },
                    {
                        label: 'Rate Limited',
                        data: ${JSON.stringify(results.map(r => r.rateLimited))},
                        backgroundColor: '#ef4444'
                    }
                ]
            },
            options: {
                responsive: true,
                scales: {
                    x: { stacked: true },
                    y: { stacked: true }
                },
                plugins: {
                    legend: { position: 'bottom' }
                }
            }
        });
    </script>
</body>
</html>
  `.trim();
}

