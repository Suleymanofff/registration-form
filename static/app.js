document.addEventListener('DOMContentLoaded', () => {
	const container = document.getElementById('container')
	const signUpBtn = document.getElementById('signUp')
	const signInBtn = document.getElementById('signIn')
	const mobileSignIn = document.getElementById('mobileSwitchToSignIn')
	const mobileSignUp = document.getElementById('mobileSwitchToSignUp')

	const togglePanel = isSignUp => {
		container.classList.toggle('right-panel-active', isSignUp)
		updateMobileSwitcher()
	}

	// Обновление мобильных ссылок
	const updateMobileSwitcher = () => {
		const isActive = container.classList.contains('right-panel-active')
		mobileSignIn.style.display = isActive ? 'inline-block' : 'none'
		mobileSignUp.style.display = isActive ? 'none' : 'inline-block'
	}

	signUpBtn.addEventListener('click', () => togglePanel(true))
	signInBtn.addEventListener('click', () => togglePanel(false))
	mobileSignIn.addEventListener('click', e => {
		e.preventDefault()
		togglePanel(false)
	})
	mobileSignUp.addEventListener('click', e => {
		e.preventDefault()
		togglePanel(true)
	})

	updateMobileSwitcher()

	const registerForm = document.getElementById('registerForm')
	const loginForm = document.getElementById('loginForm')

	// Обработчик для регистрации
	registerForm.addEventListener('submit', async e => {
		e.preventDefault()
		const formData = new FormData(registerForm)
		const data = {
			name: formData.get('name'),
			email: formData.get('email'),
			password: formData.get('password'),
		}

		try {
			const response = await fetch('/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(data),
			})

			if (response.ok) {
				alert('Регистрация прошла успешно. Теперь вы можете войти.')
				// После успешной регистрации можно переключить панель на вход
				togglePanel(false)
			} else {
				const errorText = await response.text()
				alert('Ошибка регистрации: ' + errorText)
			}
		} catch (error) {
			alert('Ошибка подключения')
		}
	})

	// Обработчик для входа
	loginForm.addEventListener('submit', async e => {
		e.preventDefault()
		const formData = new FormData(loginForm)
		const data = {
			email: formData.get('email'),
			password: formData.get('password'),
		}

		try {
			const response = await fetch('/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(data),
			})

			if (response.ok) {
				const result = await response.json()
				// Очистка страницы и вывод роли пользователя на весь экран
				document.body.innerHTML = `<div style="display:flex;align-items:center;justify-content:center;height:100vh;font-size:48px;">
					${result.role}
				</div>`
			} else {
				const errorText = await response.text()
				alert('Ошибка входа: ' + errorText)
			}
		} catch (error) {
			alert('Ошибка подключения')
		}
	})
})
