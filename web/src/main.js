import Uppy from '@uppy/core';
import DragDrop from '@uppy/drag-drop';
import StatusBar from '@uppy/status-bar';
import AwsS3Multipart from '@uppy/aws-s3-multipart';

/* CSS */

import 'tailwindcss/tailwind.css';

import '@uppy/core/dist/style.css';
import '@uppy/drag-drop/dist/style.css';
import '@uppy/status-bar/dist/style.css';

/* Components */

import Log from './log';

const log = new Log('#log-area', window.location.pathname);
document.querySelector('#log-clear').addEventListener('click', () => {
	log.clear();
});

/* Uppy */

const uppy = new Uppy({
	autoProceed: true,
});
uppy.use(DragDrop, {
	target: '#drop-area',
	height: '16rem',
});
uppy.use(StatusBar, {
	target: '#status-area',
});
uppy.use(AwsS3Multipart, {
	limit: 3,
	companionUrl: window.location.pathname,
});

/* Uppy handlers */

uppy.on('upload-success', (f, res) => {
	log.add({
		name: f.name,
		size: f.size,
		location: res.body.Location,
	});
});
