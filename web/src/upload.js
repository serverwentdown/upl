import Uppy from '@uppy/core';
import DragDrop from '@uppy/drag-drop';
import StatusBar from '@uppy/status-bar';
import AwsS3Multipart from '@uppy/aws-s3-multipart';

import Log from './log';
import Progress from './progress';

const uploadAreas = document.querySelectorAll('.upload');
uploadAreas.forEach(uploadArea => {

	/* Elements */

	const logArea = uploadArea.querySelector('.log-area');
	const logClearBtn = uploadArea.querySelector('.log-clear')
	const dropArea = uploadArea.querySelector('.drop-area');
	const statusArea = uploadArea.querySelector('.status-area');

	/* Components */

	const log = new Log(logArea, window.location.pathname, {
		empty: 'Your locally-stored file upload history is empty',
	});
	logClearBtn.addEventListener('click', () => {
		log.clear();
	});
	const progress = new Progress(statusArea);

	/* Uppy */

	const uppy = new Uppy({
		autoProceed: true,
	});
	uppy.use(DragDrop, {
		target: dropArea,
		height: '16rem',
	});
	uppy.use(AwsS3Multipart, {
		limit: 3,
		companionUrl: window.location.pathname,
	});

	uppy.on('file-added', (f, progress) => {
		console.debug(`${f.id}: Waiting...`, f)
		progress.add(f.id, {
			name: f.name,
			size: f.size,
			status: 'WAITING',
		});
	});
	uppy.on('error', (error) => {
		console.debug(`Error: ${JSON.stringify(error)}`);
	});
	uppy.on('upload-progress', (f, progress) => {
		if (!progress.uploadComplete) {
			console.debug(`${f.id}: Uploading: ${progress.percentage}%`)
		} else {
			console.debug(`${f.id}: Processing...`);
		}
		progress.update(f.id, {
			...progress,
			status: progress.uploadComplete ? 'PROCESSING' : 'UPLOADING',
		});
	});
	uppy.on('upload-error', (f, error, res) => {
		console.debug(`${f.id}: Error: ${JSON.stringify(error)} ${JSON.stringify(res)}`);
		const message = error;
		if (error.message?.includes('status: 409')) {
			message = 'A file with the same name already exists. Rename your file and try again';
		}
		progress.update(f.id, {
			error: message,
			status: 'FAILED',
		});
	});
	uppy.on('upload-retry', (id) => {
		progress.update(id, {
			status: 'RETRYING',
		});
	});
	uppy.on('upload-success', (f, res) => {
		progress.remove(f.id);
		log.add({
			name: f.name,
			size: f.size,
			location: res.body.Location,
		});
	});

});
