const updateCodeTheme = () => {
	try {
		const darkModeEnabled = document.documentElement.getAttribute('data-bs-theme') === 'dark';
		document.getElementById('light-theme').disabled = darkModeEnabled;
		document.getElementById('dark-theme').disabled = !darkModeEnabled;
	} catch {
	}
}

const getStoredTheme = () => localStorage.getItem('theme')
const setStoredTheme = theme => localStorage.setItem('theme', theme)

const getPreferredTheme = () => {
	const storedTheme = getStoredTheme()
	if (storedTheme) {
		return storedTheme
	}

	return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

const setTheme = theme => {
	if (theme === 'auto') {
		document.documentElement.setAttribute('data-bs-theme', (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))
	} else {
		document.documentElement.setAttribute('data-bs-theme', theme)
	}
	updateCodeTheme();
}

const setMermaidTheme = theme => {
	const mermaidTheme = theme === 'dark' ? 'dark' : 'default';
	mermaid.initialize({ startOnLoad: true, theme: mermaidTheme });
	mermaid.init(undefined, document.querySelectorAll('.mermaid'));
}

setTheme(getPreferredTheme())

const showActiveTheme = (theme, focus = false) => {
	const themeSwitcher = document.querySelector('#bd-theme')

	if (!themeSwitcher) {
		return
	}

	const themeSwitcherText = document.querySelector('#bd-theme-text')
	const activeThemeIcon = document.querySelector('.theme-icon-active use')
	const btnToActive = document.querySelector(`[data-bs-theme-value="${theme}"]`)
	const svgOfActiveBtn = btnToActive.querySelector('svg use').getAttribute('href')

	document.querySelectorAll('[data-bs-theme-value]').forEach(element => {
		element.classList.remove('active')
		element.setAttribute('aria-pressed', 'false')
	})

	btnToActive.classList.add('active')
	btnToActive.setAttribute('aria-pressed', 'true')
	activeThemeIcon.setAttribute('href', svgOfActiveBtn)
	const themeSwitcherLabel = `${themeSwitcherText.textContent} (${btnToActive.dataset.bsThemeValue})`
	themeSwitcher.setAttribute('aria-label', themeSwitcherLabel)

	if (focus) {
		themeSwitcher.focus()
	}
}

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
	const storedTheme = getStoredTheme()
	if (storedTheme !== 'light' && storedTheme !== 'dark') {
		setTheme(getPreferredTheme())
	}
})

function initializeDarkMode() {

	const theme = getPreferredTheme()
	setTheme(theme)
	setMermaidTheme(theme)
	showActiveTheme(theme)

	document.querySelectorAll('[data-bs-theme-value]')
		.forEach(toggle => {
			toggle.addEventListener('click', () => {
				const theme = toggle.getAttribute('data-bs-theme-value')
				setStoredTheme(theme)
				setTheme(theme)
				setMermaidTheme(theme)
				showActiveTheme(theme, true)
			})
		})
}

// Event listeners
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
	const storedTheme = getStoredTheme()
	if (storedTheme !== 'light' && storedTheme !== 'dark') {
		setTheme(getPreferredTheme())
	}
})

// Call on initial page load
document.addEventListener('DOMContentLoaded', initializeDarkMode);

// Call after HTMX content swaps
document.addEventListener('htmx:afterSettle', initializeDarkMode);
