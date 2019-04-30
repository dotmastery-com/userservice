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
                
                docker.withRegistry('https://registry-1.docker.io/v2/', 'Dockerhub') {
                    dockerImage = docker.build registry
                    dockerImage.push()
                }
            }    

         }
        }

    }
}