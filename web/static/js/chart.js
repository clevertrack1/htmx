let myChart; // Keep track of the chart instance
let loaded = false;

document.getElementById('dashboardView').addEventListener('htmx:afterSettle', (evt) => {
	if (!loaded) {
		loaded = true;

		const testChartElement = document.getElementById('testChart');
		if (!testChartElement) {
			console.error("Element with ID 'testChart' not found.");
			return;
		}

		const ctx = testChartElement.getContext('2d');
		myChart = new Chart(ctx, {
			type: 'line',
			data: {
				labels: [
					'Sunday',
					'Monday',
					'Tuesday',
					'Wednesday',
					'Thursday',
					'Friday',
					'Saturday'
				],
				datasets: [{
					data: [
						15339,
						21345,
						23489,
						12034,
						24092,
						18483,
						24003,
					],
					lineTension: 0,
					backgroundColor: 'transparent',
					borderColor: '#007bff',
					borderWidth: 4,
					pointBackgroundColor: '#007bff'
				}]
			},
			options: {
				plugins: {
					legend: {
						display: false
					},
					tooltip: {
						padding: 10 // Update padding configuration
					}
				},
				scales: {
					x: {
						grid: {
							color: '#aaa',
							lineWidth: 0.5 // Increase line width for better visibility
						}
					},
					y: {
						grid: {
							color: '#aaa',
							lineWidth: 0.5
						}
					}
				}
			}
		});
	}
});