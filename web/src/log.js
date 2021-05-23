import filesize from 'filesize';

class Log {
	constructor(target, key, messages) {
		this.key = 'log:' + key;
		this.target = target;
		this.items = [];

		const empty = document.createElement('div');
		empty.innerText = messages.empty || 'Empty';
		empty.classList.add('mt-2', 'p-2', 'bg-gray-50', 'text-gray-500', 'rounded-md', 'border', 'border-gray-400', 'flex', 'justify-center');
		this.empty = empty;

		this.localStorageLoad();
		this.render();
	}

	localStorageLoad() {
		const loaded = JSON.parse(window.localStorage.getItem(this.key) || '[]');
		this.items.push(...loaded);
	}

	localStorageSave() {
		window.localStorage.setItem(this.key, JSON.stringify(this.items));
	}

	static renderItem(item) {
		const base = document.createElement('div');	
		base.classList.add('mt-2');
		base.classList.add('flex');
		base.title = item.title;

		const url = document.createElement('input');
		url.type = 'url';
		url.value = item.location;
		url.setAttribute('readonly', '');
		url.classList.add('w-full', 'rounded-l-md', 'border-r-0', 'border-gray-400');
		url.addEventListener('click', (e) => {
			//e.target.setSelectionRange(0, e.target.value.length);
			e.target.select();
		});
		base.appendChild(url);

		const urlOpen = document.createElement('a');
		urlOpen.target = '_blank';
		urlOpen.href = item.location;
		urlOpen.classList.add('rounded-r-md', 'px-2', 'text-blue-600', 'hover:text-blue-800', 'border', 'border-l-0', 'border-gray-400', 'inline-flex', 'items-center');
		urlOpen.innerText = 'Open';
		base.appendChild(urlOpen);

		if (item.size) {
			urlOpen.classList.remove('rounded-r-md');

			const size = document.createElement('span');
			size.innerText = filesize(item.size);
			size.classList.add('text-sm', 'whitespace-nowrap', 'px-2', 'bg-gray-50', 'text-gray-500', 'rounded-r-md', 'border', 'border-gray-400', 'border-l-0', 'inline-flex', 'items-center', 'justify-center', 'w-24');
			base.appendChild(size);
		}

		return base;
	}

	render() {
		const elements = this.items.map(this.constructor.renderItem);	
		this.target.innerHTML = '';
		elements.forEach(element => {
			this.target.appendChild(element);
		});
		this.renderEmpty();
	}

	renderEmpty() {
		if (this.items.length == 0 && !this.target.contains(this.empty)) {
			this.target.appendChild(this.empty);
		} else if (this.target.contains(this.empty)) {
			this.target.removeChild(this.empty);
		}
	}

	add(item) {
		this.items.push(item);
		this.localStorageSave();
		this.target.appendChild(this.constructor.renderItem(item));
		this.renderEmpty();
	}

	clear() {
		this.items = [];
		this.localStorageSave();
		this.render();
	}
}

export default Log;
