import { PostModel } from '@/models';

class PostService {
  async getPosts(): Promise<PostModel[]> {
    const response = await fetch('https://jsonplaceholder.typicode.com/posts', {
      method: 'GET',
      next: { revalidate: 5 },
    });
    return response.json();
  }

  async getPost(postId: number): Promise<PostModel> {
    const response = await fetch(`https://jsonplaceholder.typicode.com/posts/${postId}`, {
      method: 'GET',
      next: { revalidate: 60 },
    });
    return response.json();
  }
}

export const postService = new PostService();
