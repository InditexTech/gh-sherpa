export GH_TOKEN=$(grep "extraheader" /home/runner/work/gh-sherpa/gh-sherpa/.git/config | cut -d ' ' -f 5 | cut -d ':' -f 2 | base64 -d | cut -d ':' -f 2)

gh secret list