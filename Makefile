.PHONY: check-dep
check-dep:
	cd ./awsbackend/client && ncu
	cd ./awsbackend_oauth2/client && ncu
	cd ./dexiejs-livequery && ncu
	cd ./hotupdate && ncu
	cd ./javersdemo && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./mongodb-validation && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./passwordless-dev-go/client && ncu
	cd ./pix2sketch/client && ncu
	cd ./pix2sketch/clienttalk && ncu
	cd ./pockettodo/todo && ncu
	cd ./pulumi-hetzner && ncu
	cd ./s3select/java && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./shedlockdemo && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./springai-rag && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./springai-tool && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/iac && ncu
	cd ./sqsprotobuf/go_to_java/consumer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/go_to_java_s3/consumer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/java_to_go/producer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/java_to_go_s3/producer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./transformer-js/client && ncu
	cd ./transformers-js-speech && ncu
	cd ./webpush-angular/client && ncu
	cd ./webpush-angular/server && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./xodus && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
