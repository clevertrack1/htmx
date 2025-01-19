function showDetail() {
	document.getElementById('emailListPane').classList.add('d-none');
	document.getElementById('emailUserDropdown').classList.add('d-none');
	document.getElementById('emailDetailPane').classList.remove('d-none');
	document.getElementById('emailListButton').classList.remove('d-none');
	document.getElementById('emailListButton').classList.add('d-flex');
}

function showList() {
	document.getElementById('emailListPane').classList.remove('d-none');
	document.getElementById('emailUserDropdown').classList.remove('d-none');
	document.getElementById('emailDetailPane').classList.add('d-none');
	document.getElementById('emailListButton').classList.add('d-none');
	document.getElementById('emailListButton').classList.remove('d-flex');
}

document.getElementById('emailView').addEventListener('htmx:afterSettle', (evt) => {
	//const target = evt.detail.target;
	//if (target && target.id === 'emailView' && target.dataset.initialized === 'false') {
	//	document.getElementById('emailContainer').classList.add('fade-in');
	//	target.dataset.initialized = "true";
	//}
	if (evt.detail
		&& evt.detail.pathInfo.requestPath.includes('/emails/')
		&& !evt.detail.pathInfo.requestPath.includes('search')
		&& evt.detail.requestConfig.verb === 'get'
		&& isMobile()) {
		showDetail();
	}
});

function isMobile() {
	// Bootstrap uses `d-none` combined with responsive classes (e.g., `d-md-block`)
	const testElement = document.createElement('div');
	testElement.className = 'd-none d-lg-block'; // Visible only on devices md and up
	document.body.appendChild(testElement);

	const isMobile = window.getComputedStyle(testElement).display === 'none'; // Hidden if smaller than md
	document.body.removeChild(testElement);

	return isMobile;
}

let previousScroll = 0;

// Before the swap, store the scroll position of the parent container
document.addEventListener('htmx:beforeSwap', function(e) {
	if (e.target.id === 'emailView') {
		let scrollContainer = document.getElementById('emailListPane');
		if (scrollContainer) {
			previousScroll = scrollContainer.scrollTop;
		}
	}
});

// After the swap, restore that scroll position
document.addEventListener('htmx:afterSwap', function(e) {
	if (e.target.id === 'emailView') {
		let scrollContainer = document.getElementById('emailListPane');
		scrollContainer.scrollTop = previousScroll;
	}
});
