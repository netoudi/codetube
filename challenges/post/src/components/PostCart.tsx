import Link from 'next/link';
import { Content } from '@/components/Content';
import { PostModel } from '@/models';

export function PostCard(props: { post: PostModel }) {
  return (
    <Content>
      <Link href={`/posts/${props.post.id}`}>{props.post.title}</Link>
    </Content>
  );
}
