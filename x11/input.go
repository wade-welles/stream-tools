package x11

// xdisplay is *x.Conn
// and now

//root := g.xdisplay.GetDefaultScreen().Root
//			for msg := range ch {
//				switch t := msg.(type) {
//				case *pb.InputEvent:
//					evt := msg.(*pb.InputEvent).GetKeyPressEvent()
//					key := mapping[evt.GetKey()]
//					switch evt.GetDirection() {
//					case pb.KeyPressEvent_DIRECTION_UP:
//						test.FakeInput(g.xdisplay, x.KeyReleaseEventCode, uint8(key), x.CurrentTime, root, 0, 0, 0)
//					case pb.KeyPressEvent_DIRECTION_DOWN:
//						test.FakeInput(g.xdisplay, x.KeyPressEventCode, uint8(key), x.CurrentTime, root, 0, 0, 0)
//					default:
//						log.Printf("No direction specified")
//						continue
//					}
//					g.xdisplay.Flush()
//				default:
//					log.Printf("Unexcepted type: %T", t)
//					continue
//				}
//			}
