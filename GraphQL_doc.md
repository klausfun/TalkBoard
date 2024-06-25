# GraphQL API Documentation

## Authentication

### SignUp

```graphql
mutation {
  signUp(input: { name: "test1", email: "test1@example.com", password: "password1" }) {
    id
    name
    email
  }
}
```
### SignIn
```graphql
mutation {
  signIn(input: { email: "test1@example.com", password: "password1" }) {
    token
  }
}
```

## Posts

### CreatePost
```graphql
mutation {
    createPost(input: { userId: 1, title: "post1", content: "content1", accessToComments: true }) {
        id
        title
        content
        accessToComments
    }
}
```

### GetAllPosts
```graphql
query {
    getAllPosts {
        id
        userId
        title
        content
        accessToComments
    }
}
```

## Comments
### CreateComment
```graphql
mutation {
    createComment(input: { userId: 1, postId: 2, parentCommentId: 3, content: "comment1" }) {
        id
        userId
        postId
        content
        parentCommentId
    }
}
```
### GetPostById (with Comments)
```graphql
query {
    getPostById(postId: 1, limit: 3, offset: 1) {
        post {
            id
            userId
            title
            content
            accessToComments
        }
        comments {
            id
            userId
            parentCommentId
            postId
            content
            replies {
                id
                userId
                parentCommentId
                postId
                content
                replies {
                    id
                    userId
                    parentCommentId
                    postId
                    content
                    replies {
                        id
                        userId
                        parentCommentId
                        postId
                        content
                    }
                }
            }
        }
    }
}

```