pipeline {
    agent none
    stages {
        stage('Example Build') {
            agent { docker 'maven:3.9.3-eclipse-temurin-17' }
            steps {
                echo 'Hello, Maven'
                sh 'mvn --version'
            }
        }
        stage('Example Test') {
            agent { docker 'openjdk:17-jre' }
            steps {
                echo 'Hello, JDK'
                sh 'java -version'
            }
        }
    }
}