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

	/* Components */

	const log = new Log(logArea, window.location.pathname, {
		empty: 'Your locally-stored file upload history is empty',
	});

	logClearBtn.addEventListener('click', () => {
		log.clear();
	});

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
	});
	uppy.use(AwsS3Multipart, {
		limit: 3,
		companionUrl: window.location.pathname,
	});

	uppy.on('upload-success', (f, res) => {
		log.add({
			name: f.name,
			size: f.size,
			location: res.body.Location,
		});
	});

});
