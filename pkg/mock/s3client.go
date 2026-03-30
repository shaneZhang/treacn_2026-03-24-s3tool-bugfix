package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Client 定义S3客户端接口
type S3Client interface {
	// 存储桶操作
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	DeleteBucket(ctx context.Context, params *s3.DeleteBucketInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)

	// 对象操作
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)

	// 多部分上传
	CreateMultipartUpload(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error)
	UploadPart(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error)
	ListParts(ctx context.Context, params *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error)
	CompleteMultipartUpload(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error)
	AbortMultipartUpload(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error)

	// ACL操作
	GetBucketAcl(ctx context.Context, params *s3.GetBucketAclInput, optFns ...func(*s3.Options)) (*s3.GetBucketAclOutput, error)
	PutBucketAcl(ctx context.Context, params *s3.PutBucketAclInput, optFns ...func(*s3.Options)) (*s3.PutBucketAclOutput, error)
	GetObjectAcl(ctx context.Context, params *s3.GetObjectAclInput, optFns ...func(*s3.Options)) (*s3.GetObjectAclOutput, error)
	PutObjectAcl(ctx context.Context, params *s3.PutObjectAclInput, optFns ...func(*s3.Options)) (*s3.PutObjectAclOutput, error)
}

// MockS3Client 模拟S3客户端实现
type MockS3Client struct {
	// 存储桶操作
	ListBucketsFunc       func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	CreateBucketFunc      func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	DeleteBucketFunc      func(ctx context.Context, params *s3.DeleteBucketInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
	GetBucketLocationFunc func(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
	GetBucketTaggingFunc  func(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	ListObjectsV2Func     func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObjectsFunc     func(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)

	// 对象操作
	PutObjectFunc    func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObjectFunc    func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObjectFunc func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	CopyObjectFunc   func(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
	HeadObjectFunc   func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)

	// 多部分上传
	CreateMultipartUploadFunc   func(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error)
	UploadPartFunc              func(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error)
	ListPartsFunc               func(ctx context.Context, params *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error)
	CompleteMultipartUploadFunc func(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error)
	AbortMultipartUploadFunc    func(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error)

	// ACL操作
	GetBucketAclFunc func(ctx context.Context, params *s3.GetBucketAclInput, optFns ...func(*s3.Options)) (*s3.GetBucketAclOutput, error)
	PutBucketAclFunc func(ctx context.Context, params *s3.PutBucketAclInput, optFns ...func(*s3.Options)) (*s3.PutBucketAclOutput, error)
	GetObjectAclFunc func(ctx context.Context, params *s3.GetObjectAclInput, optFns ...func(*s3.Options)) (*s3.GetObjectAclOutput, error)
	PutObjectAclFunc func(ctx context.Context, params *s3.PutObjectAclInput, optFns ...func(*s3.Options)) (*s3.PutObjectAclOutput, error)
}

// 实现S3Client接口的所有方法
func (m *MockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	if m.ListBucketsFunc != nil {
		return m.ListBucketsFunc(ctx, params, optFns...)
	}
	return &s3.ListBucketsOutput{Buckets: []types.Bucket{}}, nil
}

func (m *MockS3Client) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	if m.CreateBucketFunc != nil {
		return m.CreateBucketFunc(ctx, params, optFns...)
	}
	return &s3.CreateBucketOutput{}, nil
}

func (m *MockS3Client) DeleteBucket(ctx context.Context, params *s3.DeleteBucketInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
	if m.DeleteBucketFunc != nil {
		return m.DeleteBucketFunc(ctx, params, optFns...)
	}
	return &s3.DeleteBucketOutput{}, nil
}

func (m *MockS3Client) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	if m.GetBucketLocationFunc != nil {
		return m.GetBucketLocationFunc(ctx, params, optFns...)
	}
	return &s3.GetBucketLocationOutput{LocationConstraint: "us-east-1"}, nil
}

func (m *MockS3Client) GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	if m.GetBucketTaggingFunc != nil {
		return m.GetBucketTaggingFunc(ctx, params, optFns...)
	}
	return &s3.GetBucketTaggingOutput{TagSet: []types.Tag{}}, nil
}

func (m *MockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if m.ListObjectsV2Func != nil {
		return m.ListObjectsV2Func(ctx, params, optFns...)
	}
	return &s3.ListObjectsV2Output{Contents: []types.Object{}}, nil
}

func (m *MockS3Client) DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	if m.DeleteObjectsFunc != nil {
		return m.DeleteObjectsFunc(ctx, params, optFns...)
	}
	return &s3.DeleteObjectsOutput{}, nil
}

func (m *MockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if m.PutObjectFunc != nil {
		return m.PutObjectFunc(ctx, params, optFns...)
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if m.GetObjectFunc != nil {
		return m.GetObjectFunc(ctx, params, optFns...)
	}
	return &s3.GetObjectOutput{}, nil
}

func (m *MockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if m.DeleteObjectFunc != nil {
		return m.DeleteObjectFunc(ctx, params, optFns...)
	}
	return &s3.DeleteObjectOutput{}, nil
}

func (m *MockS3Client) CopyObject(ctx context.Context, params *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error) {
	if m.CopyObjectFunc != nil {
		return m.CopyObjectFunc(ctx, params, optFns...)
	}
	return &s3.CopyObjectOutput{}, nil
}

func (m *MockS3Client) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if m.HeadObjectFunc != nil {
		return m.HeadObjectFunc(ctx, params, optFns...)
	}
	return &s3.HeadObjectOutput{}, nil
}

func (m *MockS3Client) CreateMultipartUpload(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	if m.CreateMultipartUploadFunc != nil {
		return m.CreateMultipartUploadFunc(ctx, params, optFns...)
	}
	return &s3.CreateMultipartUploadOutput{UploadId: aws.String("test-upload-id")}, nil
}

func (m *MockS3Client) UploadPart(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	if m.UploadPartFunc != nil {
		return m.UploadPartFunc(ctx, params, optFns...)
	}
	return &s3.UploadPartOutput{ETag: aws.String("test-etag")}, nil
}

func (m *MockS3Client) ListParts(ctx context.Context, params *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
	if m.ListPartsFunc != nil {
		return m.ListPartsFunc(ctx, params, optFns...)
	}
	return &s3.ListPartsOutput{Parts: []types.Part{}}, nil
}

func (m *MockS3Client) CompleteMultipartUpload(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	if m.CompleteMultipartUploadFunc != nil {
		return m.CompleteMultipartUploadFunc(ctx, params, optFns...)
	}
	return &s3.CompleteMultipartUploadOutput{}, nil
}

func (m *MockS3Client) AbortMultipartUpload(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	if m.AbortMultipartUploadFunc != nil {
		return m.AbortMultipartUploadFunc(ctx, params, optFns...)
	}
	return &s3.AbortMultipartUploadOutput{}, nil
}

func (m *MockS3Client) GetBucketAcl(ctx context.Context, params *s3.GetBucketAclInput, optFns ...func(*s3.Options)) (*s3.GetBucketAclOutput, error) {
	if m.GetBucketAclFunc != nil {
		return m.GetBucketAclFunc(ctx, params, optFns...)
	}
	return &s3.GetBucketAclOutput{
		Owner: &types.Owner{
			DisplayName: aws.String("test-owner"),
			ID:          aws.String("test-id"),
		},
		Grants: []types.Grant{},
	}, nil
}

func (m *MockS3Client) PutBucketAcl(ctx context.Context, params *s3.PutBucketAclInput, optFns ...func(*s3.Options)) (*s3.PutBucketAclOutput, error) {
	if m.PutBucketAclFunc != nil {
		return m.PutBucketAclFunc(ctx, params, optFns...)
	}
	return &s3.PutBucketAclOutput{}, nil
}

func (m *MockS3Client) GetObjectAcl(ctx context.Context, params *s3.GetObjectAclInput, optFns ...func(*s3.Options)) (*s3.GetObjectAclOutput, error) {
	if m.GetObjectAclFunc != nil {
		return m.GetObjectAclFunc(ctx, params, optFns...)
	}
	return &s3.GetObjectAclOutput{
		Owner: &types.Owner{
			DisplayName: aws.String("test-owner"),
			ID:          aws.String("test-id"),
		},
		Grants: []types.Grant{},
	}, nil
}

func (m *MockS3Client) PutObjectAcl(ctx context.Context, params *s3.PutObjectAclInput, optFns ...func(*s3.Options)) (*s3.PutObjectAclOutput, error) {
	if m.PutObjectAclFunc != nil {
		return m.PutObjectAclFunc(ctx, params, optFns...)
	}
	return &s3.PutObjectAclOutput{}, nil
}
