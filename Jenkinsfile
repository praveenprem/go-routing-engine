def NODE = 'Slave'
def SSH_KEY_ID = '34eff5cc-5c4c-46ab-a541-600c84729870'

pipeline {
    agent {
        node {
            label NODE
        }
    }

    options {
        timeout(time: 15, unit: 'MINUTES')
        buildDiscarder(logRotator(numToKeepStr: '20'))
        disableConcurrentBuilds(abortPrevious: true)
        timestamps()
    }

    parameters {
        string(name: 'version', defaultValue: '', description: 'Release version (X.X.X).\n0.0.1')
    }

    stages {

        stage("Pre-build validation") {
            steps {
                script {
                    def match = params.version =~ /^[0-9]+\.[0-9]+\.[0-9]+$/
                    if (!match.find()) {
                        error("Release version doesn't match Semantic Versioning.")
                    }
                }
            }
        }

        stage('Start Release') {
            steps {
                print 'Release start'
                sh "git checkout -B release/${params.version}"
            }
        }

        stage('Update CHANGELOG') {
            steps {
                sh """
                sed -i -e 's/\\(##.*\\)\\(unreleased\\)\\(.*\\)\$/\\1\\2\\3\\n\\n\\1v${params.version}\\3/g' \\
                    -e 's/\\(\\[.*\\]\\)\\(.*\\/\\)\\(.*\\)\\.\\.\\.\\(.*\\)\$/\\1\\2v${params.version}...HEAD\\n\\[v${params.version}\\]\\2\\3...v${params.version}/g' CHANGELOG.md
                """
            }
        }

        stage('Release finish') {
            steps {
                sh "git add CHANGELOG.md"
                sh "git commit -S -m \"Update CHANGELOG and version\""
                sh "git checkout master"
                sh "git merge -S --no-ff --no-commit release/${params.version}"
                sh "git commit -S -m \"Merge release/${params.version}\""
                sh "git branch -D release/${params.version}"
                sh "git tag -s v${params.version} -m \"Release ${params.version}\""
                sh "git checkout develop"
                sh "git merge -S --no-commit master"
                print 'Release finish'
            }
        }

        stage('Push to origin') {
            input {
                message 'Push to origin?'
                ok 'Yes'
            }
            steps {
                sshagent(credentials: [SSH_KEY_ID]) {
                    sh 'git push origin master'
                    sh 'git push origin develop'
                    sh "git push origin v${params.version}"
                }
            }
        }
    }

    post {
        success {
            slackSend color: "good", message: "${env.JOB_BASE_NAME} - #${env.BUILD_NUMBER} Success - (<${env.BUILD_URL}|Open>)"
        }

        aborted {
            slackSend color: "warning", message: "${env.JOB_BASE_NAME} - #${env.BUILD_NUMBER} aborted by user - (<${env.BUILD_URL}|Open>)"
        }

        failure {
            slackSend color: "danger", message: "${env.JOB_BASE_NAME} - #${env.BUILD_NUMBER} Failure - (<${env.BUILD_URL}|Open>)"
        }

        cleanup {
            deleteDir()
        }
    }
}
