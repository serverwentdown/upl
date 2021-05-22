import Uppy from '@uppy/core';
import DragDrop from '@uppy/drag-drop';
import StatusBar from '@uppy/status-bar';
import AwsS3Multipart from '@uppy/aws-s3-multipart';

import '@uppy/core/dist/style.css';
import '@uppy/drag-drop/dist/style.css';
import '@uppy/status-bar/dist/style.css';
import './main.css';

const uppy = new Uppy({
	autoProceed: true,
});
uppy.use(DragDrop, {
	target: '#drop-area',
});
uppy.use(StatusBar, {
	target: '#status-area',
});
uppy.use(AwsS3Multipart, {
	limit: 3,
	companionUrl: '.',
});

uppy.on('upload-success', (f, res) => {
	console.log(f, res);
});
