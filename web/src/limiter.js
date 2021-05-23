class Limiter {
	constructor(delay=500) {
		this.delay = delay;

		this.interval = null;
		this.next = null;
	}

	call(fn) {
		if (this.interval == null) {
			fn();
			this.enableClock();
		} else {
			this.next = fn;
		}
	}

	checkCall() {
		if (this.next == null) {
			return this.disableClock();
		}

		// Function wanted to be called
		this.next();
		this.next = null;
	}

	enableClock() {
		this.interval = setInterval(() => {
			this.checkCall();
		}, this.delay);
	}

	disableClock() {
		clearInterval(this.interval);
		this.interval = null;
	}
}

export default Limiter;
