import Log from './log';

async function throwForStatus(res) {
	if (!res.ok) {
		const text = await res.text();
		throw new Error(text);
	}
	return res;
}

const createAreas = document.querySelectorAll('.create');
createAreas.forEach(createArea => {

	/* Elements */

	const logArea = createArea.querySelector('.log-area');
	const logClearBtn = createArea.querySelector('.log-clear');
	const noticeArea = createArea.querySelector('.notice-area');

	/* Components */

	const log = new Log(logArea, window.location.pathname, {
		empty: 'Your locally-stored dropbox creation history is empty',
	});

	logClearBtn.addEventListener('click', () => {
		log.clear();
	});

	/* Form */

	function showError(error='') {
		let message = error.message || error.toString();
		if (message !== '') {
			noticeArea.classList.remove('hidden');
			noticeArea.innerText = message;
		} else {
			noticeArea.classList.add('hidden');
			noticeArea.innerText = message;
		}
	}

	createArea.addEventListener('submit', e => {
		e.preventDefault();

		// Clear previous errors
		showError();

		const data = new FormData(e.target);
		fetch(e.target.action, {
			method: e.target.method,
			body: data,
		})
			.then(throwForStatus)
			.then(res => res.text())
			.then(id => {
				log.add({
					location: window.location.origin + '/' + id,
					title: `
Endpoint: ${data.get('Endpoint')}
Region: ${data.get('Region')}
Access key: ${data.get('AccessKey')}
Canned ACL: ${data.get('ACL')}}
Prefix: ${data.get('Prefix')}
Expires: ${data.get('ExpiresNumber')}${data.get('ExpiresUnits')}
		`.trim(),
				});
			})
			.catch(showError);
	});

});
