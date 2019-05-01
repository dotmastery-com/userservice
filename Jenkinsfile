pipeline {
    agent none



    environment {
        registry = "dotmastery/userservice"
    }
    stages {


        stage('Build and Push Docker Image') {
          agent any  

          steps {
            script {
                
              //  docker.withRegistry('https://registry-1.docker.io/v2/', 'Dockerhub') {
              //      dockerImage = docker.build registry
               //     dockerImage.push()
               // }

                sh 'whoami'
                sh 'cd /app/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit'    
                sh 'pwd'
                sh '/app/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit/oc login https://api.starter-us-east-1.openshift.com --token=U0F4Fy17A5TNfTHviU4NNQYiifLzIfnW9YpovIfDMG8'
                sh '/app/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit/oc import-image myuserservice:latest --from=dotmastery/userservice --confirm'
            }   


         }
        }

    }
}