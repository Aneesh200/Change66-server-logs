#!/usr/bin/env node

/**
 * Setup Verification Script
 * Run this script to verify your logs-ui setup is correct
 */

const fs = require('fs');
const path = require('path');

console.log('🔍 Checking Logs UI Setup...\n');

let hasErrors = false;

// Check 1: Node.js version
console.log('1. Checking Node.js version...');
const nodeVersion = process.version;
const majorVersion = parseInt(nodeVersion.slice(1).split('.')[0]);
if (majorVersion >= 18) {
  console.log(`   ✓ Node.js ${nodeVersion} (minimum: 18.x)`);
} else {
  console.log(`   ✗ Node.js ${nodeVersion} is too old. Please upgrade to 18.x or later.`);
  hasErrors = true;
}
console.log();

// Check 2: package.json exists
console.log('2. Checking package.json...');
if (fs.existsSync('package.json')) {
  console.log('   ✓ package.json found');
} else {
  console.log('   ✗ package.json not found');
  hasErrors = true;
}
console.log();

// Check 3: node_modules exists
console.log('3. Checking dependencies...');
if (fs.existsSync('node_modules')) {
  console.log('   ✓ node_modules directory found');
  
  // Check for key dependencies
  const keyDeps = ['next', 'react', 'react-dom', 'tailwindcss'];
  let allDepsInstalled = true;
  
  keyDeps.forEach(dep => {
    if (!fs.existsSync(path.join('node_modules', dep))) {
      console.log(`   ✗ Missing dependency: ${dep}`);
      allDepsInstalled = false;
      hasErrors = true;
    }
  });
  
  if (allDepsInstalled) {
    console.log('   ✓ All key dependencies installed');
  }
} else {
  console.log('   ✗ node_modules not found. Run: npm install');
  hasErrors = true;
}
console.log();

// Check 4: .env.local configuration
console.log('4. Checking environment configuration...');
if (fs.existsSync('.env.local')) {
  console.log('   ✓ .env.local file found');
  
  const envContent = fs.readFileSync('.env.local', 'utf8');
  
  if (envContent.includes('NEXT_PUBLIC_API_URL')) {
    const match = envContent.match(/NEXT_PUBLIC_API_URL=(.+)/);
    if (match && match[1].trim()) {
      console.log(`   ✓ API URL configured: ${match[1].trim()}`);
    } else {
      console.log('   ⚠ NEXT_PUBLIC_API_URL is empty');
    }
  } else {
    console.log('   ✗ NEXT_PUBLIC_API_URL not found in .env.local');
    hasErrors = true;
  }
  
  if (envContent.includes('NEXT_PUBLIC_API_KEY')) {
    console.log('   ✓ API key configured');
  } else {
    console.log('   ℹ API key not configured (optional if server doesn\'t require auth)');
  }
} else {
  console.log('   ✗ .env.local file not found');
  console.log('   → Create it with: echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local');
  hasErrors = true;
}
console.log();

// Check 5: Required directories
console.log('5. Checking project structure...');
const requiredDirs = ['app', 'components', 'lib'];
let allDirsExist = true;

requiredDirs.forEach(dir => {
  if (fs.existsSync(dir)) {
    console.log(`   ✓ ${dir}/ directory found`);
  } else {
    console.log(`   ✗ ${dir}/ directory not found`);
    allDirsExist = false;
    hasErrors = true;
  }
});
console.log();

// Check 6: Key files
console.log('6. Checking key files...');
const requiredFiles = [
  'app/page.tsx',
  'app/layout.tsx',
  'components/logs-table.tsx',
  'components/logs-filter.tsx',
  'components/pagination.tsx',
  'lib/api.ts',
  'lib/types.ts',
];

let allFilesExist = true;

requiredFiles.forEach(file => {
  if (fs.existsSync(file)) {
    console.log(`   ✓ ${file}`);
  } else {
    console.log(`   ✗ ${file} not found`);
    allFilesExist = false;
    hasErrors = true;
  }
});
console.log();

// Final summary
console.log('═══════════════════════════════════════');
if (!hasErrors) {
  console.log('✅ All checks passed! Your setup is ready.');
  console.log('\nNext steps:');
  console.log('  1. Make sure your log server is running');
  console.log('  2. Enable CORS on the server (ENABLE_CORS=true)');
  console.log('  3. Run: npm run dev');
  console.log('  4. Open: http://localhost:3000');
} else {
  console.log('❌ Some checks failed. Please fix the issues above.');
  console.log('\nQuick fixes:');
  console.log('  • Run: npm install');
  console.log('  • Create .env.local with your API URL');
  console.log('  • Make sure you\'re in the logs-ui directory');
}
console.log('═══════════════════════════════════════\n');

process.exit(hasErrors ? 1 : 0);


