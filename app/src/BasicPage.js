import Header from "./Header";

export default function BasicPage(props) {
  return (
    <div>
      <Header signOut={props.signOut}/>

      {props.children}
    </div>
  );
}
