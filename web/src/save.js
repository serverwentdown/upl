import Limiter from './limiter';

class Saver {
	constructor(target, inputs, key='') {
		this.target = target;
		this.inputs = inputs;
		this.key = 'save:' + key;

		this.limiter = new Limiter();

		this.updateState = this.updateState.bind(this);
		this.input = this.input.bind(this);
		this.save = this.save.bind(this);

		this.loadState();
		this.target.addEventListener('input', this.updateState);
	}

	loadState() {
		if (window.localStorage.getItem(this.key) !== null) {
			this.target.checked = true;
		} else {
			this.target.checked = false;
		}
		this.updateBind();
	}

	updateState() {
		if (this.target.checked) {
			this.save();
		} else {
			this.clear();
		}
		this.updateBind();
	}

	updateBind() {
		if (this.target.checked) {
			this.inputs.forEach(input => input.addEventListener('input', this.input));
		} else {
			this.inputs.forEach(input => input.removeEventListener('input', this.input));
		}
	}

	input() {
		this.limiter.call(this.save);
	}

	load() {
		const values = JSON.parse(window.localStorage.getItem(this.key) || '{}');
		for (const input of this.inputs) {
			if (input.name in values) {
				input.value = values[input.name];
			}
		}
	}

	save() {
		const values = {};
		for (const input of this.inputs) {
			values[input.name] = input.value;
		}
		window.localStorage.setItem(this.key, JSON.stringify(values));
	}

	clear() {
		window.localStorage.removeItem(this.key);
	}
}

const saveInputs = document.querySelectorAll('[data-save]');
saveInputs.forEach(saveInput => {

	const inputNames = saveInput.dataset.save.split(',');

	const inputs = inputNames.map(inputName => document.querySelector(`[name="${inputName}"]`));
	const saver = new Saver(saveInput, inputs, saveInput.dataset.save);
	saver.load();

});
