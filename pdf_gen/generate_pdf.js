const puppeteer = require('puppeteer');
const path = require('path');

(async () => {
  console.log('Launching browser...');
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();

  const htmlPath = 'file://' + path.resolve(__dirname, 'explanation.html');
  console.log('Loading HTML from:', htmlPath);

  await page.goto(htmlPath, { waitUntil: 'networkidle0', timeout: 60000 });

  // Wait a bit for fonts/styles to settle
  await new Promise(r => setTimeout(r, 2000));

  const outputPath = path.resolve(__dirname, '../UserManagement_Code_Explanation.pdf');
  console.log('Generating PDF at:', outputPath);

  await page.pdf({
    path: outputPath,
    format: 'A4',
    printBackground: true,
    margin: {
      top: '20mm',
      bottom: '20mm',
      left: '15mm',
      right: '15mm'
    },
    displayHeaderFooter: true,
    headerTemplate: '<div style="font-size:9px;color:#6366f1;width:100%;text-align:center;font-family:sans-serif;padding-top:5px">User Management API — Complete Code Explanation</div>',
    footerTemplate: '<div style="font-size:9px;color:#94a3b8;width:100%;text-align:center;font-family:sans-serif;padding-bottom:5px">Page <span class="pageNumber"></span> of <span class="totalPages"></span></div>'
  });

  await browser.close();
  console.log('✅ PDF generated successfully!');
  console.log('📄 Saved to: ' + outputPath);
})();
