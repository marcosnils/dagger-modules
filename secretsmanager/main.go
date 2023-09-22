package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awssm "github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Secretsmanager struct {
	Key    string
	Secret string
}

func (m *Secretsmanager) Auth(key, secret string) *Secretsmanager {
	m.Key = key
	m.Secret = secret

	return m
}

func (m *Secretsmanager) GetSecret(name string) (string, error) {
	config := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(m.Key, m.Secret, ""),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return "", err
	}
	svc := awssm.New(sess)

	input := &awssm.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	value, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}
	return *(value.SecretString), nil
}

func (m *Secretsmanager) PutSecret(name, value string) (*Secretsmanager, error) {
	config := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(m.Key, m.Secret, ""),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	svc := awssm.New(sess)

	input := &awssm.PutSecretValueInput{
		SecretId:     aws.String(name),
		SecretString: aws.String(value),
	}

	_, err = svc.PutSecretValue(input)
	return m, err
}

func (c *Container) WithAWSSecret(key, secret, name, envName string) (*Container, error) {
	sm := &Secretsmanager{
		Key:    key,
		Secret: secret,
	}

	s, err := sm.GetSecret(name)
	if err != nil {
		return nil, err
	}

	dagSecret := dag.SetSecret(envName, s)
	return c.WithSecretVariable(envName, dagSecret), nil
}
