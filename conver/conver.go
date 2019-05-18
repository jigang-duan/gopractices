package conver

type Converter interface {
	Reverse(entitys ...interface{}) (model interface{})
	Convert(models ...interface{}) (entity interface{})
}

