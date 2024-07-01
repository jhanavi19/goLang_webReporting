
document.getElementById('reportForm').addEventListener('submit', function(event) {
    event.preventDefault();
    generateReport();
});

async function generateReport() {
    const formData = {
        name: document.getElementById('name').value,
        occupation: document.getElementById('occupation').value,
        hobby: document.getElementById('hobby').value
    };

    try {
        const response = await fetch('http://127.0.0.1:8080/generate-report', {
            method: 'POST',
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            throw new Error('Failed to generate report');
        }

        
        fetchAllReports();
    } catch (error) {
        console.error('Error generating report:', error);
        document.getElementById('reportResult').innerHTML = '<p>Failed to generate report.</p>';
    }
}

async function fetchAllReports() {
    try {
        const response = await fetch('http://localhost:8080/api/reports');
        if (!response.ok) {
            throw new Error('Failed to fetch reports');
        }
        const allReports = await response.json();
        console.log('Fetched Reports:', allReports);
        displayAllReports(allReports);
    } catch (error) {
        console.error('Error fetching reports:', error);
        document.getElementById('reportResult').innerHTML = '<p>Failed to fetch reports.</p>';
    }
}

function displayAllReports(reports) {
    const reportResult = document.getElementById('reportResult');
    let tableHtml = '<h2>All Reports</h2><table border="1"><tr><th>Name</th><th>Occupation</th><th>Hobby</th></tr>';

    reports.forEach(report => {
        tableHtml += `<tr><td>${report.name}</td><td>${report.occupation}</td><td>${report.hobby}</td></tr>`;
    });

    tableHtml += '</table>';
    tableHtml += '<button onclick="downloadAsPDF()">Download as PDF</button>';
    reportResult.innerHTML = tableHtml;
}

async function downloadAsPDF() {
    try {
        const doc = new jsPDF();
        const columns = ['Name', 'Occupation', 'Hobby'];
        const rows = [];

        document.querySelectorAll('table tr').forEach((row, index) => {
            if (index > 0) {
                const rowData = [];
               
                row.querySelectorAll('td').forEach(cell => {
                    rowData.push(cell.textContent.trim());
                });
             
                rows.push(rowData);
            }
        });

        console.log('Extracted Rows:', rows);
        if (rows.length > 0) {
            doc.autoTable({
                head: [columns],
                body: rows,
            });
            doc.save('reports.pdf');
        } else {
            throw new Error('No data found in the table to generate PDF.');
        }
    } catch (error) {
        console.error('Error generating PDF:', error);
        alert('Failed to generate PDF. Please try again.');
    }
}


