#!/usr/bin/env groovy
import java.text.SimpleDateFormat
@Library('custom-lib@feat_cn_methods') _


def args = [:]

/*
Arguments need to be modified:
- stack_name                         # stack name that the service is located in rancher
- image_name                         # path to your image in registry
- rancher_access_key_credential_id   # credential id of rancher access key
- rancher_secret_key_credential_id   # credential id of rancher secret key
- email_to                           # email receivers
- dockerfile_path                    # path to the Dockerfile, '.' by default
*/

args.stack_name = 'message-hub-backend'
args.image_name = 'message-hub/message-hub-mock'
args.service_name = scm.getUserRemoteConfigs()[0].getUrl().tokenize('/').last().split("\\.")[0]
args.github_project_id = 'dktunited/' + args.service_name

args.rancher_access_key_credential_id = env.BRANCH_NAME == 'master' ? 'messagehub_rancher_access_key_production' : 'messagehub_rancher_access_key_preprod'
args.rancher_secret_key_credential_id = env.BRANCH_NAME == 'master' ? 'messagehub_rancher_secret_key_production' : 'messagehub_rancher_secret_key_preprod'

args.email_to   = 'alex.lu@decathlon.com'

args.dockerfile_path = '.'

/*
    Basically, the following args do NOT need to be modified
    change the proxy_host to Azure proxy server if the ci/cd platform (NOT your application) is running on Azure
*/

args.image_version = ''

args.rancher_cli = '/opt/rancher_cli/rancher'
args.rancher_compose_cli = '/opt/rancher_cli/rancher-compose'

args.sonar_project_key = 'dktunited_' + scm.getUserRemoteConfigs()[0].getUrl().tokenize('/').last().split("\\.")[0]

args.proxy_host = 'proxy-internet-aws-china-production.subsidia.org'
args.no_proxy_host = 'localhost|127.0.0.1|*.subsidia.org'
args.registry_local = 'registry-cn-local.subsidia.org'
args.registry = 'registry-cn.subsidia.org'

pipeline{

    agent any

    tools {
        maven 'maven3.6'
        jdk 'java11'
    }

    options {
        // Basic options
        timestamps()
        disableConcurrentBuilds()
        buildDiscarder(logRotator(numToKeepStr:'5'))
    }

    stages {

        //stage('Unit Test and Quality Analysis') {
        //  parallel {
        //      stage('Unit Test') {
        //          steps {
        //              script {
        //                    sh 'mvn org.jacoco:jacoco-maven-plugin:prepare-agent test -fae'
        //                }
        //           }
        //            post {
        //                always {
        //                    jacoco execPattern: '**/target/coverage-reports/jacoco-ut.exec'
        //                }
        //            }
        //        }
        //        stage('Sonar - Preview For Commit') {
        //            when {
        //                expression { env.CHANGE_ID != null }
        //            }
        //            steps {
        //                withCredentials([string(credentialsId: 'sonar_login', variable: 'sonar_login')]) {
        //                    script {
        //                        withSonarQubeEnv('sonar-cloud') {
        //                            def commit_id = sh(script: "git rev-parse refs/remotes/origin/${env.BRANCH_NAME}", returnStdout: true).trim()
        //                            log.info "git commit id = $commit_id"
        //                            sh "mvn sonar:sonar -Dsonar.coverage.jacoco.xmlReportPaths='target/site/jacoco-ut/jacoco.xml' -Ddockerfile.skip -Dsonar.host.url=https://sonarcloud.io -Dsonar.projectKey=$args.sonar_project_key -Dsonar.login=$sonar_login -Dsonar.organization=dktunited -Dsonar.pullrequest.branch=${env.BRANCH_NAME} -Dsonar.pullrequest.key=${env.CHANGE_ID} -Dsonar.scm.revision=$commit_id -Dsonar.pullrequest.provider=github -Dhttp.proxyPort=3128 -Dhttp.proxyHost=$args.proxy_host -Dhttps.proxyPort=3128 -Dhttps.proxyHost=$args.proxy_host -Dhttp.nonProxyHosts='$args.no_proxy_host'"
        //                        }
        //                    }
        //                }
        //            }
        //        }
        //        stage('Sonar - Full Analysis on Develop or Master Branch') {
        //            when {
        //                anyOf {
        //                    branch 'master'
        //                    branch 'develop'
        //                }
        //            }
        //            steps {
        //                withCredentials([string(credentialsId: 'sonar_login', variable: 'sonar_login')]) {
        //                    script {
        //                        withSonarQubeEnv('sonar-cloud') {
        //                            sh "mvn sonar:sonar -Dsonar.coverage.jacoco.xmlReportPaths='target/site/jacoco-ut/jacoco.xml' -Ddockerfile.skip -Dsonar.host.url=https://sonarcloud.io -Dsonar.projectKey=$args.sonar_project_key -Dsonar.login=$sonar_login -Dsonar.organization=dktunited -Dsonar.branch.name=${env.BRANCH_NAME} -Dhttp.proxyPort=3128 -Dhttp.proxyHost=$args.proxy_host -Dhttps.proxyPort=3128 -Dhttps.proxyHost=$args.proxy_host -Dhttp.nonProxyHosts='$args.no_proxy_host'"
        //                        }
        //                    }
        //                }
        //            }
        //        }
        //    }
        //}

        stage('Build Docker Image and Prepare Rancher Deployment') {
            parallel {
                stage('Build docker image') {
                    when {
                        anyOf {
                            branch 'master'
                            branch 'develop'
                        }
                    }
                    steps {
                        script {

                            args.image_version = version.getBuildImageVersion()
                            log.info "version number : $args.image_version"

                            image = docker.build("$args.image_name:$args.image_version", "--build-arg http_proxy=http://proxy-internet-aws-china-production.subsidia.org:3128 --build-arg https_proxy=http://proxy-internet-aws-china-production.subsidia.org:3128 $args.dockerfile_path")

                            docker.withRegistry("https://$args.registry_local",'nexusAccount'){
                                image.push("$args.image_version")
                            }

                            sh "docker rmi -f $args.image_name:$args.image_version"
                            sh "docker rmi -f $args.registry_local/$args.image_name:$args.image_version"

                            def imageLine = "$args.registry/$args.image_name:$args.image_version $env.WORKSPACE/$args.dockerfile_path/Dockerfile"
                            writeFile file: 'anchore_images', text: imageLine
                            anchore name: 'anchore_images'
                        }
                    }
                }

                stage('Preprare Rancher Deployment') {
                    when {
                        anyOf {
                            branch 'master'
                            branch 'develop'
                        }
                    }
                    steps{
                        script {
                            try {
                                sh "$args.rancher_cli --version && $args.rancher_compose_cli --version"
                            } catch (Exception err) {
                                rancher.prepare()
                                args.rancher_cli = './rancher'
                                args.rancher_compose_cli = './rancher-compose'
                            }
                            def stacks = rancher.getStacks(args)
                            if (rancher.isStackExists(stacks, args.stack_name) == false) rancher.createStack(args)
                        }
                    }
                }
            }
        }

        stage('Download Compose and Delete Unhealthy Services') {
            parallel {
                stage('Download the compose') {
                    when {
                        anyOf {
                            branch 'master'
                            branch 'develop'
                        }
                    }
                    steps{
                        script {
                            rancher.exportCompose(args)
                        }
                    }
                }

                stage('Delete unhealthy services'){
                    when {
                        anyOf {
                            branch 'master'
                            branch 'develop'
                        }
                    }
                    steps {
                        script {
                            rancher.deleteUnhealthyService(args)
                        }
                    }
                }
            }
        }

        stage('Deploy') {
            when {
                anyOf {
                    branch 'master'
                    branch 'develop'
                }
            }
            steps{
                script {
                    rancher.deploy(args)
                }
            }
        }
    }

    post {
        success {
            script {
                email.send_notification(args, "Successful")
            }
        }

        failure {
            script {
                email.send_notification(args, "Failed")
            }
        }

        always {
            cleanWs()
        }
    }
}
