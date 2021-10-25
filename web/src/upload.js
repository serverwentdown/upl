import Uppy from '@uppy/core';
import DragDrop from '@uppy/drag-drop';
import StatusBar from '@uppy/status-bar';
import AwsS3Multipart from '@uppy/aws-s3-multipart';

import Log from './log';

const uploadAreas = document.querySelectorAll('.upload');
uploadAreas.forEach(uploadArea => {

	/* Elements */

	const logArea = uploadArea.querySelector('.log-area');
	const logClearBtn = uploadArea.querySelector('.log-clear')
	const dropArea = uploadArea.querySelector('.drop-area');
	const statusArea = uploadArea.querySelector('.status-area');
	const noticeArea = uploadArea.querySelector('.notice-area');

	/* Components */

	const log = new Log(logArea, window.location.pathname, {
		empty: 'Your locally-stored file upload history is empty',
	});

	logClearBtn.addEventListener('click', () => {
		log.clear();
	});

	/* Error */

	function showError(error='') {
		let message = error.message || error.toString();
		if (message !== '') {
			noticeArea.classList.remove('hidden');
			noticeArea.innerText = message;
		} else {
			noticeArea.classList.add('hidden');
			noticeArea.innerText = message;
		}
		window.scrollTo({ top: 0 });
	}

	/* Uppy */

	const uppy = new Uppy({
		autoProceed: true,
	});
	uppy.use(DragDrop, {
		target: dropArea,
		height: '16rem',
	});
	uppy.use(StatusBar, {
		target: statusArea,
		showProgressDetails: true,
	});
	uppy.use(AwsS3Multipart, {
		limit: 3,
		companionUrl: window.location.pathname,
	});

	uppy.on('file-added', (f, progress) => {
		console.debug(`${f.id}: Waiting...`)
	});
	uppy.on('upload-progress', (f, progress) => {
		if (!progress.uploadComplete) {
			console.debug(`${f.id}: Uploading: ${progress.percentage}%`)
		} else {
			console.debug(`${f.id}: Processing...`);
		}
	});
	uppy.on('upload-error', (f, error, res) => {
		console.debug(`${f.id}: Error: ${error}`);
		window.e = { f, error, res };
		if (error.message?.includes('status: 409')) {
			showError('A file with the same name already exists. Rename your file and try again');
			return;
		}
		showError(error);
	});
	uppy.on('upload-retry', (id) => {
		showError();
	});
	uppy.on('upload-success', (f, res) => {
		showError();
		log.add({
			name: f.name,
			size: f.size,
			location: res.body.Location,
		});
	});

});
