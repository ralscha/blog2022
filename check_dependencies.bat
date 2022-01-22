cd D:\ws\blog2022\sqsprotobuf\go_to_java\consumer
call checkdep

cd D:\ws\blog2022\sqsprotobuf\go_to_java\producer
call make upgrade-libraries

cd D:\ws\blog2022\sqsprotobuf\go_to_java_s3\consumer
call checkdep

cd D:\ws\blog2022\sqsprotobuf\go_to_java_s3\producer
call make upgrade-libraries


cd D:\ws\blog2022\sqsprotobuf\iac
call ncu

cd D:\ws\blog2022\sqsprotobuf\java_to_go\consumer
call make upgrade-libraries

cd D:\ws\blog2022\sqsprotobuf\java_to_go\producer
call checkdep

cd D:\ws\blog2022\sqsprotobuf\java_to_go_s3\consumer
call make upgrade-libraries

cd D:\ws\blog2022\sqsprotobuf\java_to_go_s3\producer
call checkdep

cd D:\ws\blog2022\dexiejs-livequery
call ncu

call D:\ws\blog2022\s3select\go
call make upgrade-libraries

call D:\ws\blog2022\s3select\java
call checkdep
