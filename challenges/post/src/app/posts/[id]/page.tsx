import { Container } from '@/components/Container';
import { Content } from '@/components/Content';
import { Heading } from '@/components/Heading';
import { postService } from '@/post-service';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export default async function PostDetails(props: any) {
  const { id } = await props.params;
  const post = await postService.getPost(+id);

  return (
    <Container>
      <Heading>{post.title}</Heading>
      <Content>
        <p className="text-lg">{post.body}</p>
      </Content>
    </Container>
  );
}
