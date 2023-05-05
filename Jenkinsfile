def NODE = 'Slave'
def GPG_KEY_ID = '53298CC2E45BABE9AAE5FC342EB97A2D4CD98917'
def GPG_PASSPHRASE_ID = '76854055-a980-4292-a1a1-e9476da9a946'

pipeline {
    agent {
        node {
            label NODE
        }
    }

    parameters {
        string(name: 'version', defaultValue: '', description: 'Release version (X.X.X).\n0.0.1')
    }

    stages {

        stage('GPG Config') {
            environment {
                gpg_passphrase = credentials(GPG_PASSPHRASE_ID)
            }
            steps {
                sh "/usr/lib/gnupg2/gpg-preset-passphrase --preset --passphrase ${gpg_passphrase} $GPG_KEY_ID"
            }
        }

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
                sh "git checkout -B release/${params.version}"
            }
        }

        stage('Update CHANGELOG') {
            steps {
                sh """
                sed -i -e 's/\\(##.*\\)\\(unreleased\\)\\(.*\\)\$/\\1\\2\\3\\n\\n\\1${params.version}\\3/g' \\
                    -e 's/\\(\\[.*\\]\\)\\(.*\\/\\)\\(.*\\)\\.\\.\\.\\(.*\\)\$/\\1\\21.0.0...HEAD\\n\\[1.0.0\\]\\2\\3...${params.version}/g' CHANGELOG.md
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
                sh "git tag -s ${params.version}"
                sh "git checkout develop"
                sh "git merge -S --no-commit master"
            }
        }
    }

    post {
        cleanup {
            deleteDir()
        }
    }
}
