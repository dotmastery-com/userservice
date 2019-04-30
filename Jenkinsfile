node {
    
    // Install the desired Go version
    echo "Installing go"
    def root = tool name: '1.12.1', type: 'go'
    
    try{
       
        
        ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/UserService") {
            withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}", "GOROOT=${root}","PATH+GO=${root}/bin"]) {
                env.PATH="${GOPATH}/bin:$PATH"
                
                stage('Checkout'){
                    echo 'Checking out SCM'
                    checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[url: 'https://github.com/dotmastery-com/userservice.git']]])
                }
                
                stage('Pre Test'){
                    echo 'Pulling Dependencies'
            
//                    sh 'go version'
 //                   sh 'go get github.com/rs/cors'
 //                   sh 'go get github.com/gorilla/mux'
 //                   sh 'go get github.com/rs/xid'
 //                   sh 'go get gopkg.in/mgo.v2'
 //                   sh 'go get gopkg.in/mgo.v2/bson'
                
                 //   sh 'go get -u github.com/golang/dep/cmd/dep'
                 //   sh 'dep init && dep ensure'

                }
        
                stage('Test'){
                    
                    //List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
                    //Push our project files relative to ./src
                  //  sh 'cd $GOPATH && go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org > projectPaths'
                    
                    //Print them with 'awk '$0="./src/"$0' projectPaths' in order to get full relative path to $GOPATH
                  //  def paths = sh returnStdout: true, script: """awk '\$0="./src/"\$0' projectPaths"""
                    
                  //  echo 'Vetting'

                  //  sh """cd $GOPATH && go tool vet ${paths}"""

                  //  echo 'Linting'
                  //  sh """cd $GOPATH && golint ${paths}"""
                    
                  // echo 'Testing'
                  //  sh """cd $GOPATH && go test -race -cover ${paths}"""
                }
            
                stage('Build and Push Docker Image') {

                    echo 'step 1'

                    agent any  

                    environment {
                        registry = "dotmastery/userservice"
                    }

echo 'step 2'
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
    }catch (e) {
        // If there was an exception thrown, the build failed
        currentBuild.result = "FAILED"
        
        
    } finally {
        
    }
}