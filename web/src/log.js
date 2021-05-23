import filesize from 'filesize';

class Log {
	constructor(selector, key) {
		this.key = key;
		this.target = document.querySelector(selector);
		this.items = [];

		this.localStorageLoad();
		this.render();
	}

	localStorageLoad() {
		const loaded = JSON.parse(window.localStorage.getItem('log' + this.key) || '[]');
		this.items.push(...loaded);
	}

	localStorageSave() {
		window.localStorage.setItem('log' + this.key, JSON.stringify(this.items));
	}

	static renderItem(item) {
		const base = document.createElement('div');	
		base.classList.add('log-item');

		const url = document.createElement('input');
		url.value = item.location;
		url.setAttribute('readonly', '');
		url.classList.add('log-url');
		url.addEventListener('click', (e) => {
			e.target.setSelectionRange(0, e.target.value.length);
		});
		base.appendChild(url);

		const size = document.createElement('span');
		size.innerText = filesize(item.size);
		size.classList.add('log-size');
		base.appendChild(size);

		return base;
	}

	render() {
		const elements = this.items.map(this.constructor.renderItem);	
		this.target.innerHTML = '';
		elements.forEach(element => {
			this.target.appendChild(element);
		});
	}

	add(item) {
		this.items.push(item);
		this.localStorageSave();
		this.target.appendChild(this.constructor.renderItem(item));
	}

	clear() {
		this.items = [];
		this.localStorageSave();
		this.render();
	}
}

export default Log;
