FROM scratch

ADD deployer /bin/deployer
ADD deployerd /bin/deployerd

CMD ["deployerd"]
