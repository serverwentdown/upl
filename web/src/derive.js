import * as time from './time';

class InvalidInputError extends Error {
	constructor(problems) {
		super(`Invalid input: ${problems.join(', ')}`);
		this.problems = problems;
	}
}

class Deriver {
	constructor(output, inputs, notice=null) {
		this.output = output;
		this.inputs = inputs;
		this.notice = notice;
	}

	get values() {
		return this.inputs.map(input => input.value);
	}

	bind() {
		this.inputs.forEach(input => {
			input.addEventListener('input', () => {
				this.update();
			});
		});
	}

	update() {
		try {
			const output = this.derive();
			this.output.value = output;
			this.showError();
		} catch (e) {
			console.error(e);
			this.showError(e);
		}
	}

	showError(error='') {
		if (!this.notice) {
			return;
		}
		let message = error.message || error.toString();
		if (Array.isArray(error.problems)) {
			message = 'Invalid input: ' + error.problems.join(', ');
		}
		this.notice.innerText = message;
	}
}

class DurationDeriver extends Deriver {
	derive() {
		const [ number, units ] = this.values;
		const n = parseInt(number, 10);
		if (!isFinite(n)) {
			throw new InvalidInputError([`provided duration is not a number`]);
		}
		switch (units) {
		case 's':
			return n * time.SECOND;
		case 'm':
			return n * time.MINUTE;
		case 'h':
			return n * time.HOUR;
		case 'd':
			return n * 24 * time.HOUR;
		}
		throw new InvalidInputError([`unit ${units} is not valid`]);
	}
}

const derivers = {
	duration: DurationDeriver,
};

const deriveInputs = document.querySelectorAll('[data-derive]');
deriveInputs.forEach(deriveInput => {

	const [ type, ...inputNames ] = deriveInput.dataset.derive.split(',');
	if (!type in derivers) {
		return;
	}

	const inputs = inputNames.map(inputName => document.querySelector(`[name="${inputName}"]`));
	const notice = document.querySelector(deriveInput.dataset.deriveNotice);

	const deriver = new derivers[type](deriveInput, inputs, notice);
	deriver.bind();
	deriver.update();

});
