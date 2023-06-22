FROM mcr.microsoft.com/windows/nanoserver:ltsc2022
COPY k8s-wait-for-multi.exe /
ENTRYPOINT ["/k8s-wait-for-multi.exe"]