.PHONY: all
all: upgrade check-dep

.PHONY: upgrade
upgrade:
	make -C ./sqsprotobuf/go_to_java_s3/producer upgrade-libraries
	make -C ./sqsprotobuf/java_to_go/consumer upgrade-libraries
	make -C ./sqsprotobuf/java_to_go_s3/consumer upgrade-libraries
	make -C ./sqsprotobuf/go_to_java/producer upgrade-libraries
	make -C ./s3select/go upgrade-libraries
	make -C ./awsbackend/iac upgrade-libraries
	make -C ./awsbackend/lambda upgrade-libraries
	make -C ./awsbackend_oauth2/iac upgrade-libraries
	make -C ./awsbackend_oauth2/lambda upgrade-libraries
	make -C ./hibp-go/api_server upgrade-libraries
	make -C ./hibp-go/bloom upgrade-libraries
	make -C ./hibp-go/pebble upgrade-libraries

.PHONY: check-dep
check-dep:
	cd ./s3select/java && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/iac && ncu
	cd ./dexiejs-livequery && ncu
	cd ./sqsprotobuf/go_to_java/consumer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/go_to_java_s3/consumer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/java_to_go/producer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/java_to_go_s3/producer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./awsbackend/client && ncu
	cd ./awsbackend_oauth2/client && ncu
	cd ./hotupdate && ncu
	cd ./passwordless-dev-go/client && ncu
	cd ./dexiejs-livequery && ncu
	cd ./transformers-js-speech && ncu
	cd ./webpush-angular/client && ncu
	cd ./webpush-angular/server && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./transformer-js/client && ncu
	cd ./pix2sketch/clienttalk && ncu
	cd ./pix2sketch/client && ncu
	cd ./transformer-js/client && ncu
	cd ./transformers-js-speech && ncu	
	cd ./webpush-angular/client && ncu
	cd ./webpush-angular/server && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./springai-tool && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates	