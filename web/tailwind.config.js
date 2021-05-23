module.exports = {
	purge: [
		'./*.tmpl',
		'./src/*.js',
	],
	plugins: [
		require('@tailwindcss/forms'),
	],
};
