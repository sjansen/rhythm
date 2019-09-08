sync:
	aws --profile buscaluz s3 cp --recursive webui/build/static/  s3://rhythm.tuntap.net-media/static/
	for I in webui/build/*.*; do aws --profile buscaluz s3 cp $$I s3://rhythm.tuntap.net-media/; done
