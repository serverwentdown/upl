import filesize from 'filesize';

class Log {
	constructor(target, key) {
		this.key = key;
		this.target = target;
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
		base.classList.add('mt-1');
		base.classList.add('flex');

		const url = document.createElement('input');
		url.type = 'url';
		url.value = item.location;
		url.setAttribute('readonly', '');
		url.classList.add('w-full');
		url.classList.add('rounded-l-md');
		url.classList.add('border-gray-400');
		url.addEventListener('click', (e) => {
			e.target.setSelectionRange(0, e.target.value.length);
		});
		base.appendChild(url);

		const size = document.createElement('span');
		size.innerText = filesize(item.size);
		size.classList.add('text-sm');
		size.classList.add('whitespace-nowrap');
		size.classList.add('px-2');
		size.classList.add('bg-gray-50');
		size.classList.add('text-gray-500');
		size.classList.add('rounded-r-md');
		size.classList.add('border');
		size.classList.add('border-gray-400');
		size.classList.add('border-l-0');
		size.classList.add('inline-flex');
		size.classList.add('items-center');
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
