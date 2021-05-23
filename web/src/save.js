import Limiter from './limiter';

class Saver {
	constructor(target, inputs, key='') {
		this.target = target;
		this.inputs = inputs;
		this.key = '';

		this.limiter = new Limiter();

		this.updateState = this.updateState.bind(this);
		this.input = this.input.bind(this);
		this.save = this.save.bind(this);

		this.target.addEventListener('input', this.updateState);
		this.updateState();
	}

	updateState() {
		if (this.target.checked) {
			this.inputs.forEach(input => input.addEventListener('input', this.input));
			this.save();
		} else {
			this.inputs.forEach(input => input.removeEventListener('input', this.input));
			this.clear();
		}
	}

	input() {
		this.limiter.call(this.save);
	}

	load() {
		const values = JSON.parse(window.localStorage.getItem('save' + this.key) || '{}');
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
		window.localStorage.setItem('save' + this.key, JSON.stringify(values));
	}

	clear() {
		window.localStorage.removeItem('save' + this.key);
	}
}

const saveInputs = document.querySelectorAll('[data-save]');
saveInputs.forEach(saveInput => {

	const inputNames = saveInput.dataset.save.split(',');

	const inputs = inputNames.map(inputName => document.querySelector(`[name="${inputName}"]`));
	const saver = new Saver(saveInput, inputs);
	saver.load();

});
