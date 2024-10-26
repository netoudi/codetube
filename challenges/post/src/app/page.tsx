import { Container } from '@/components/Container';
import { Heading } from '@/components/Heading';
import { PostCard } from '@/components/PostCart';
import { postService } from '@/post-service';

export default async function Home() {
  const posts = await postService.getPosts();

  return (
    <Container>
      <Heading>Posts</Heading>
      {posts.map((post) => {
        return <PostCard key={post.id} post={post} />;
      })}
    </Container>
  );
}
