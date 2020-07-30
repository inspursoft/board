from setuptools import setup, find_packages

setup(
  name="gitlab-helper",
  version="0.1.0",
  packages=find_packages(),
  install_requires=[
    "requests", 
    "paramiko"
  ],
)