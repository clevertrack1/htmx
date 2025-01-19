(function() {

	document.addEventListener('DOMContentLoaded', function() {

		const searchModal = document.getElementById('searchModal');
		console.log(searchModal);
		searchModal.addEventListener('shown.bs.modal', function() {
			var searchInput = document.getElementById('searchInput');
			searchInput.focus();
			searchInput.select();
		})

		const bsSearchModal = new bootstrap.Modal(searchModal);
		document.addEventListener('keydown', function(event) {
			if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
				event.preventDefault();
				bsSearchModal.show();
			}
		});

	});

})();
