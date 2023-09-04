import avatar from "gradient-avatar";

const Avatar = ({ id, size = 32 }: { id: string; size?: number }) => {
  const avatarSVG = avatar(id);
  const dataUri = `data:image/svg+xml,${encodeURIComponent(avatarSVG)}`;
  return (
    <div >
      <img class=" rounded-full" src={dataUri} alt="Avatar" width={size} height={size} />
    </div>
  );
};

export default Avatar;
