Project_Name = 'TWLotteryCrawler'
Iris = 'node.rayer.idv.tw'
Iris_OCR1 = 'node1.rayer.idv.tw'
Iris_OCR2 = 'node2.rayer.idv.tw'
pipeline {
    agent any
    parameters {
        string defaultValue: 'master', description: 'Branch name', name: 'branch', trim: true
        string defaultValue: 'cli.app', description: 'CLI app name', name: 'cli_app_name', trim: true
        string defaultValue: 'demoApp.app', description: 'UI app name', name: 'demo_app_name', trim: true
    }

    stages {
        stage('Fetch from github') {
            steps {
                slackSend message: "Project ${Project_Name} start to be built"
                git credentialsId: '26c5c0a0-d02d-4d77-af28-761ffb97c5cc', url: 'https://github.com/Rayer/TWLotteryCrawler.git', branch: "${params.branch}"
            }
        }
        stage('Unit test') {
            steps {
                sh label: 'go version', script: 'go version'
                sh label: 'install gocover-cobertura', script: 'go get github.com/t-yuki/gocover-cobertura'
                sh label: 'go unit test', script: 'go test --coverprofile=cover.out'
                sh label: 'convert coverage xml', script: '~/go/bin/gocover-cobertura < cover.out > coverage.xml'
            }
        }
        stage ("Extract test and coverage results") {
            steps {
                cobertura coberturaReportFile: 'coverage.xml'
            }
        }

        stage('build and archive executable') {
            steps {
                sh label: 'show version', script: 'go version'
                sh label: 'build cli', script: "go build  -o bin/${params.cli_app_name} ./lotteryCli/*.go"
                //sh label: 'build app', script: "go build -o bin/${params.demo_app_name} ./lotteryUi/*.go"
                archiveArtifacts artifacts: 'bin/*', fingerprint: true, followSymlinks: false, onlyIfSuccessful: true
            }
        }
    }

   post {
        aborted {
            slackSend message: "Project ${Project_Name} aborted."
        }
        success {
            slackSend message: "Project ${Project_Name} is built successfully."
        }
        failure {
            slackSend message: "Project ${Project_Name} is failed to be built."
        }
    }
}