FROM mcr.microsoft.com/windows/nanoserver:ltsc2019
COPY k8s-wait-for-multi.exe /
ENTRYPOINT ["/k8s-wait-for-multi.exe"]