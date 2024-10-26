export function Container(props: { children: React.ReactNode }) {
  return <main className="flex flex-col gap-8 w-[1200px] py-8 mx-auto">{props.children}</main>;
}
