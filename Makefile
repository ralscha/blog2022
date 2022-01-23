.PHONY: all
all: upgrade check-dep

.PHONY: upgrade
upgrade:
	make -C ./sqsprotobuf/go_to_java_s3/producer upgrade-libraries
	make -C ./sqsprotobuf/java_to_go/consumer upgrade-libraries
	make -C ./sqsprotobuf/java_to_go_s3/consumer upgrade-libraries
	make -C ./sqsprotobuf/go_to_java/producer upgrade-libraries
	make -C ./s3select/go upgrade-libraries

.PHONY: check-dep
check-dep:
	cd ./s3select/java && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/iac && ncu
	cd ./dexiejs-livequery && ncu
	cd ./sqsprotobuf/go_to_java/consumer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/go_to_java_s3/consumer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/java_to_go/producer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
	cd ./sqsprotobuf/java_to_go_s3/producer && ./mvnw.cmd versions:display-dependency-updates && ./mvnw.cmd versions:display-plugin-updates
