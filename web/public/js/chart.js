let loaded = false;
document.getElementById('dashboardView').addEventListener('htmx:afterSettle', (evt) => {
	if (!loaded) {
		loaded = true;
		const ctx = document.getElementById('testChart').getContext('2d');
		const myChart = new Chart(ctx, {
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
						boxPadding: 3
					}
				},
				scales: {
					x: {
						grid: {
							color: '#aaa', // Change grid color here
							lineWidth: 0.25
						}
					},
					y: {
						grid: {
							color: '#aaa', // Change grid color here
							lineWidth: 0.25
						}
					}
				}
			}
		});

	}
})
