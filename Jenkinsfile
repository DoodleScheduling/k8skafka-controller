import java.util.regex.Pattern
import org.jenkinsci.plugins.pipeline.modeldefinition.Utils

podTemplate(label: 'k8skafka-controller',
  containers: [
    containerTemplate(
      name: 'golang',
      image: 'bitnami/golang:1.15',
      ttyEnabled: true
    ),
    containerTemplate(
      name: 'docker',
      image: 'docker:latest',
      ttyEnabled: true
    ),
    containerTemplate(
      name: 'helm',
      command: '/bin/ash',
      image: 'alpine/helm:latest',
      ttyEnabled: true
    ),
  ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
  ]
) {
  node ('k8skafka-controller') {
    ansiColor("xterm") {
      stage('checkout') {
        checkout(scm)
      
        container('docker') {
          dockerAuth()
        }
      }

      stage("build") {
        container('golang') {
          sh 'make all'
        }

        container('helm') {
          sh 'helm lint chart/k8skafka-controller'
        }
      }

      stage("publish") {
        if (!env.TAG_NAME) {
          echo "skip packaging for no tagged release"
        } else {
          def (_,major,minor,patch,group,label,build) = (env.TAG_NAME =~ /^v(\d{1,3})\.(\d{1,3})\.(\d{1,3})(?:(-([A-Za-z0-9]+)))?$/)[0]

          if (!major && !minor && !patch) {
            throw new Exception("Invalid tag detected, requires semantic version")
          }

          version = "$major.$minor.$patch$group"

          container('docker') {
            sh "docker build . -t nexus.doodle.com:5000/devops/k8skafka-controller:${env.TAG_NAME}"
            sh "docker push nexus.doodle.com:5000/devops/k8skafka-controller:${env.TAG_NAME}"
          }

          container('helm') {
            bumpChartVersion(version)
            bumpImageVersion(env.TAG_NAME)

            tgz="k8skafka-controller-${version}.tgz"
            sh "mkdir chart/k8skafka-controller/crds"
            sh "cp config/crd/bases/* chart/k8skafka-controller/crds"
            sh "helm package chart/k8skafka-controller"

          }

          container('golang') {
            if (label) {
              publish(tgz, "helm-staging")
            } else {
              publish(tgz, "helm-staging")
              publish(tgz, "helm-production")
            }
          }
        }
      }
    }
  }
}

void dockerAuth() {
  // nexus repository
  withCredentials([[
                       $class          : 'UsernamePasswordMultiBinding',
                       credentialsId   : 'nexus',
                       usernameVariable: 'NEXUS_USER',
                       passwordVariable: 'NEXUS_PASSWORD'
                   ]]) {
    sh "docker login nexus.doodle.com:5000 -u ${env.NEXUS_USER} -p ${env.NEXUS_PASSWORD}"
  }

  // docker hub
  withCredentials([[
                       $class          : 'UsernamePasswordMultiBinding',
                       credentialsId   : 'dockerhub',
                       usernameVariable: 'DOCKERHUB_USER',
                       passwordVariable: 'DOCKERHUB_PASSWORD'
                   ]]) {
    sh "docker login -u ${env.DOCKERHUB_USER} -p ${env.DOCKERHUB_PASSWORD}"
  }
}

def bumpImageVersion(String version) {
  echo "Update image tag"
  def valuesFile = "./chart/k8skafka-controller/values.yaml"
  def valuesData = readYaml file: valuesFile
  valuesData.image.tag = version

  sh "rm $valuesFile"
  writeYaml file: valuesFile, data: valuesData
}

def bumpChartVersion(String version) {
  // Bump chart version
  echo "Update chart version"
  def chartFile = "./chart/k8skafka-controller/Chart.yaml"
  def chartData = readYaml file: chartFile
  chartData.version = version
  chartData.appVersion = version

  sh "rm $chartFile"
  writeYaml file: chartFile, data: chartData
}

def publish(String tgz, String repository) {
  echo "Push chart ${tgz} to helm repository ${repository}"

  withCredentials([[
    $class          : 'UsernamePasswordMultiBinding',
    credentialsId   : 'nexus',
    usernameVariable: 'NEXUS_USER',
    passwordVariable: 'NEXUS_PASSWORD'
  ]]) {
    sh "curl -u \"${env.NEXUS_USER}:${env.NEXUS_PASSWORD}\" https://nexus.doodle.com/repository/${repository}/ --upload-file $tgz --fail"
  }
}